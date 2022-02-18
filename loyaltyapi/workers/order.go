package workers

import (
	"context"
	"errors"
	"github.com/putalexey/gophermart/loyaltyapi/models"
	"github.com/putalexey/gophermart/loyaltyapi/repository"
	"github.com/putalexey/gophermart/loyaltyapi/services"
	"go.uber.org/zap"
	"sync"
	"time"
)

type OrderWorker struct {
	ctx            context.Context
	logger         *zap.SugaredLogger
	repo           repository.Repositorier
	accrualService *services.Accrual
	jobTTL         time.Duration
}

func New(
	ctx context.Context,
	logger *zap.SugaredLogger,
	repo repository.Repositorier,
	jobTTL time.Duration,
	accrualService *services.Accrual,
) *OrderWorker {
	return &OrderWorker{
		ctx:            ctx,
		logger:         logger,
		repo:           repo,
		jobTTL:         jobTTL,
		accrualService: accrualService,
	}
}

func (o *OrderWorker) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		o.logger.Info("order worker started")

		t := time.NewTicker(5 * time.Second)
		defer t.Stop()

		for {
			o.ProcessJob()
			select {
			case <-t.C:
			case <-o.ctx.Done():
				return
			}
		}

	}()

	wg.Wait()
	o.logger.Info("order worker stopped")
}

func (o *OrderWorker) ProcessJob() {
	job, err := o.repo.TakeJob(o.ctx, o.jobTTL)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			o.logger.Info("no jobs available")
			return
		}
		o.logger.Error(err)
		return
	}
	o.logger.Infof("processing job %s", job.UUID)
	defer func() {
		o.logger.Infof("processed job %s", job.UUID)
	}()

	order, err := o.repo.GetOrder(o.ctx, job.OrderUUID)
	if err != nil {
		o.logger.Error(err)
		return
	}

	reqCtx, cancel := context.WithTimeout(o.ctx, 3*time.Second)
	defer cancel()

	orderStatus, err := o.accrualService.GetOrderStatus(reqCtx, order.Number)
	if err != nil {
		tooManyRequests := services.ErrTooManyRequests{}
		if errors.As(err, &tooManyRequests) {
			o.logger.Warn(err)
			if tooManyRequests.RetryAfter != nil {
				job.ProceedAt = *tooManyRequests.RetryAfter
				err := o.repo.UpdateJob(o.ctx, job)
				if err != nil {
					o.logger.Error(err)
					return
				}
			}
			return
		}
		if errors.Is(err, services.ErrNotFound) {
			job.Tries = 0
			err = o.repo.UpdateJob(o.ctx, job)
			if err != nil {
				o.logger.Error(err)
			}
			return
		}
		o.logger.Error(err)
		return
	}

	isFinalStatus := false
	doDeposit := false
	switch orderStatus.Status {
	case "REGISTERED":
	case "INVALID": // — заказ не принят к расчёту, и вознаграждение не будет начислено;
		isFinalStatus = true
		order.Status = models.OrderStatusInvalid
	case "PROCESSING": // — расчёт начисления в процессе;
		order.Status = models.OrderStatusProcessing
	case "PROCESSED": // — расчёт начисления окончен;
		isFinalStatus = true
		if order.Status != models.OrderStatusProcessed {
			doDeposit = true
			order.Status = models.OrderStatusProcessed
			order.Accrual = orderStatus.Accrual
		}
	}

	_, err = o.repo.SaveOrder(o.ctx, order)
	if err != nil {
		o.logger.Error(err)
		return
	}

	if doDeposit {
		deposit := &models.Deposit{
			UserUUID: order.UserUUID,
			Sum:      order.Accrual,
		}
		_, err = o.repo.BalanceDeposit(o.ctx, deposit)
		if err != nil {
			o.logger.Error(err)
			return
		}
	}

	if isFinalStatus {
		err = o.repo.DeleteJob(o.ctx, job)
		if err != nil {
			o.logger.Error(err)
			return
		}
		return
	}
	job.Tries = 0
	err = o.repo.UpdateJob(o.ctx, job)
	if err != nil {
		o.logger.Error(err)
	}
}
