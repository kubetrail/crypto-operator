package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cryptov1beta1 "github.com/kubetrail/crypto-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *AccountReconciler) FinalizeStatus(ctx context.Context, clientObject client.Object) error {
	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	object, ok := clientObject.(*cryptov1beta1.Account)
	if !ok {
		err := fmt.Errorf("cientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	// Update the status of the object if not terminating
	if object.Status.Phase != phaseTerminating {
		object.Status = cryptov1beta1.AccountStatus{
			Phase:      phaseTerminating,
			Conditions: object.Status.Conditions,
			Message:    "object is marked for deletion",
			Reason:     reasonObjectMarkedForDeletion,
		}
		if err := r.Status().Update(ctx, object); err != nil {
			reqLogger.Error(err, "failed to update object status")
			return err
		} else {
			reqLogger.Info("updated object status")
			return ObjectUpdated
		}
	}

	return nil
}

func (r *AccountReconciler) FinalizeResources(ctx context.Context, clientObject client.Object, req ctrl.Request) error {
	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	object, ok := clientObject.(*cryptov1beta1.Account)
	if !ok {
		err := fmt.Errorf("clientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	reqLogger.Info("object deleted")

	var found bool
	// Update the status of the object if pending
	for i, condition := range object.Status.Conditions {
		if condition.Reason == reasonDeleted {
			object.Status.Conditions[i].LastTransitionTime = metav1.Time{Time: time.Now()}
			found = true
			break
		}
	}

	if !found {
		condition := metav1.Condition{
			Type:               conditionTypeRuntime,
			Status:             metav1.ConditionTrue,
			ObservedGeneration: 0,
			LastTransitionTime: metav1.Time{Time: time.Now()},
			Reason:             reasonDeleted,
			Message:            "deleted object",
		}
		object.Status = cryptov1beta1.AccountStatus{
			Phase:      object.Status.Phase,
			Conditions: append(object.Status.Conditions, condition),
			Message:    "deleted object",
			Reason:     reasonDeleted,
		}

		if err := r.Status().Update(ctx, object); err != nil {
			reqLogger.Error(err, "failed to update object status")
			return err
		} else {
			reqLogger.Info("updated object status")
			return ObjectUpdated
		}
	}

	return nil
}

func (r *AccountReconciler) RemoveFinalizer(ctx context.Context, clientObject client.Object) error {
	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	controllerutil.RemoveFinalizer(clientObject, finalizer)
	if err := r.Update(ctx, clientObject); err != nil {
		reqLogger.Error(err, "failed to remove finalizer")
		return err
	}
	reqLogger.Info("finalizer removed")
	return ObjectUpdated
}

func (r *AccountReconciler) AddFinalizer(ctx context.Context, clientObject client.Object) error {
	if controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	controllerutil.AddFinalizer(clientObject, finalizer)
	if err := r.Update(ctx, clientObject); err != nil {
		reqLogger.Error(err, "failed to add finalizer")
		return err
	}
	reqLogger.Info("finalizer added")
	return ObjectUpdated
}

func (r *AccountReconciler) InitializeStatus(ctx context.Context, clientObject client.Object) error {
	reqLogger := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		err := fmt.Errorf("finalizer not found")
		reqLogger.Error(err, "failed to detect finalizer")
		return err
	}

	object, ok := clientObject.(*cryptov1beta1.Account)
	if !ok {
		err := fmt.Errorf("cientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	// Update the status of the object if none exists
	found := false
	for _, condition := range object.Status.Conditions {
		if condition.Reason == reasonFinalizerAdded {
			found = true
			break
		}
	}

	if !found {
		object.Status = cryptov1beta1.AccountStatus{
			Phase: phasePending,
			Conditions: []metav1.Condition{
				{
					Type:               conditionTypeObject,
					Status:             metav1.ConditionTrue,
					ObservedGeneration: 0,
					LastTransitionTime: metav1.Time{Time: time.Now()},
					Reason:             reasonFinalizerAdded,
					Message:            "object initialized",
				},
			},
			Message: "object initialized",
			Reason:  reasonObjectInitialized,
		}
		if err := r.Status().Update(ctx, object); err != nil {
			reqLogger.Error(err, "failed to update object status")
			return err
		} else {
			reqLogger.Info("updated object status")
			return ObjectUpdated
		}
	}

	return nil
}

func (r *AccountReconciler) ReconcileResources(ctx context.Context, clientObject client.Object, req ctrl.Request) error {
	reqLogger := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		err := fmt.Errorf("finalizer not found")
		reqLogger.Error(err, "failed to detect finalizer")
		return err
	}

	object, ok := clientObject.(*cryptov1beta1.Account)
	if !ok {
		err := fmt.Errorf("cientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	// if object has label key crypto.kubetrail.io/group then
	// filter coins with that particular label key: value pair
	var coinLabel string
	for k, v := range object.Labels {
		if k == label {
			coinLabel = v
			break
		}
	}

	coins := &cryptov1beta1.CoinList{}
	if len(coinLabel) > 0 {
		if err := r.List(
			ctx, coins,
			client.InNamespace(object.Namespace),
			client.MatchingLabels{
				label: coinLabel,
			},
		); err != nil {
			reqLogger.Error(err, "failed to get list of coins")
			return err
		}
	} else {
		if err := r.List(
			ctx, coins,
			client.InNamespace(object.Namespace),
		); err != nil {
			reqLogger.Error(err, "failed to get list of coins")
			return err
		}
	}

	var balance float64
	var numCoins int
	coinList := make([]string, 0, len(coins.Items))
	for _, coin := range coins.Items {
		if len(coin.Status.Meta.Balance) == 0 {
			continue
		}

		if coin.Status.Phase != phaseRunning {
			continue
		}

		b, err := strconv.ParseFloat(coin.Status.Meta.Balance, 64)
		if err != nil {
			reqLogger.Error(err, "failed to parse coin balance")
			return err
		}

		balance += b
		numCoins++
		coinList = append(coinList, coin.Name)
	}

	found := false
	foundIndex := -1
	updated := false
	for i, condition := range object.Status.Conditions {
		if condition.Type == conditionTypeRuntime &&
			condition.Status == metav1.ConditionTrue &&
			condition.Reason == reasonSynced {
			found = true
			if time.Since(condition.LastTransitionTime.Time) > time.Minute {
				foundIndex = i
			}
			break
		}
	}

	condition := metav1.Condition{
		Type:               conditionTypeRuntime,
		Status:             metav1.ConditionTrue,
		ObservedGeneration: 0,
		LastTransitionTime: metav1.Time{Time: time.Now()},
		Reason:             reasonSynced,
		Message:            "synced coin prices",
	}

	if !found {
		conditions := object.Status.Conditions
		object.Status = cryptov1beta1.AccountStatus{
			Meta: cryptov1beta1.AccountStatusMeta{
				NumCoins: numCoins,
				CoinList: coinList,
				Balance:  fmt.Sprintf("%f", balance),
			},
			Phase:      phaseRunning,
			Conditions: append(conditions, condition),
			Message:    "synced coin prices",
			Reason:     reasonSynced,
		}
		updated = true
	} else {
		if foundIndex >= 0 {
			object.Status.Conditions[foundIndex] = condition
			conditions := object.Status.Conditions
			object.Status = cryptov1beta1.AccountStatus{
				Meta: cryptov1beta1.AccountStatusMeta{
					NumCoins: numCoins,
					CoinList: coinList,
					Balance:  fmt.Sprintf("%f", balance),
				},
				Phase:      phaseRunning,
				Conditions: conditions,
				Message:    "synced coin prices",
				Reason:     reasonSynced,
			}
			updated = true
		}
	}

	if updated {
		if err := r.Status().Update(ctx, object); err != nil {
			reqLogger.Error(err, "failed to update object status")
			return err
		} else {
			rateLimit(
				fmt.Sprintf("%s-%s", object.Name, object.Namespace),
				time.Hour*24,
				func() {
					reqLogger.Info("updated object status")
				},
			)
			return ObjectUpdated
		}
	}

	return nil
}
