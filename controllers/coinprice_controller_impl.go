package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	cryptov1beta1 "github.com/kubetrail/crypto-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *CoinReconciler) FinalizeStatus(ctx context.Context, clientObject client.Object) error {
	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	object, ok := clientObject.(*cryptov1beta1.Coin)
	if !ok {
		err := fmt.Errorf("cientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	// Update the status of the object if not terminating
	if object.Status.Phase != phaseTerminating {
		object.Status = cryptov1beta1.CoinStatus{
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

func (r *CoinReconciler) FinalizeResources(ctx context.Context, clientObject client.Object, req ctrl.Request) error {
	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		return nil
	}

	reqLogger := log.FromContext(ctx)

	object, ok := clientObject.(*cryptov1beta1.Coin)
	if !ok {
		err := fmt.Errorf("clientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	reqLogger.Info("object deleted")

	var found bool
	// Update the status of the object if pending
	for i, condition := range object.Status.Conditions {
		if condition.Reason == reasonDeletedCoin {
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
			Reason:             reasonDeletedCoin,
			Message:            "deleted coin",
		}
		object.Status = cryptov1beta1.CoinStatus{
			Phase:      object.Status.Phase,
			Conditions: append(object.Status.Conditions, condition),
			Message:    "deleted coin",
			Reason:     reasonDeletedCoin,
		}
	}

	if err := r.Status().Update(ctx, object); err != nil {
		reqLogger.Error(err, "failed to update object status")
		return err
	} else {
		reqLogger.Info("updated object status")
		return ObjectUpdated
	}
}

func (r *CoinReconciler) RemoveFinalizer(ctx context.Context, clientObject client.Object) error {
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

func (r *CoinReconciler) AddFinalizer(ctx context.Context, clientObject client.Object) error {
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

func (r *CoinReconciler) InitializeStatus(ctx context.Context, clientObject client.Object) error {
	reqLogger := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		err := fmt.Errorf("finalizer not found")
		reqLogger.Error(err, "failed to detect finalizer")
		return err
	}

	object, ok := clientObject.(*cryptov1beta1.Coin)
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
		object.Status = cryptov1beta1.CoinStatus{
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

func (r *CoinReconciler) ReconcileResources(ctx context.Context, clientObject client.Object, req ctrl.Request) error {
	reqLogger := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(clientObject, finalizer) {
		err := fmt.Errorf("finalizer not found")
		reqLogger.Error(err, "failed to detect finalizer")
		return err
	}

	object, ok := clientObject.(*cryptov1beta1.Coin)
	if !ok {
		err := fmt.Errorf("cientObject to object type assertion error")
		reqLogger.Error(err, "failed to get object instance")
		return err
	}

	url := fmt.Sprintf(
		"https://api.coinbase.com/v2/prices/%s-%s/spot",
		strings.ToLower(object.Spec.Ticker),
		strings.ToLower(object.Spec.Currency),
	)
	response, err := http.Get(url)
	if err != nil {
		msg := "failed to get coin price"
		err := fmt.Errorf("%s: %w", msg, err)
		reqLogger.Error(err, msg)
		return err
	}

	jb, err := io.ReadAll(response.Body)
	if err != nil {
		msg := "failed to read coin price json response"
		err := fmt.Errorf("%s: %w", msg, err)
		reqLogger.Error(err, msg)
		return err
	}

	price := &coinprice{}
	if err := json.Unmarshal(jb, price); err != nil {
		msg := "failed to parse json response from coin price query"
		err := fmt.Errorf("%s-%w", msg, err)
		reqLogger.Error(err, msg)
		return err
	}

	n, err := strconv.ParseFloat(object.Spec.NumCoins, 64)
	if err != nil {
		msg := "failed to parse object numCoins"
		err := fmt.Errorf("%s: %s", msg, err)
		reqLogger.Error(err, msg)
		return err
	}

	p, err := strconv.ParseFloat(price.Data.Amount, 64)
	if err != nil {
		msg := "failed to parse price data amount"
		err := fmt.Errorf("%s: %s", msg, err)
		reqLogger.Error(err, msg)
		return err
	}

	found := false
	foundIndex := -1
	updated := false
	for i, condition := range object.Status.Conditions {
		if condition.Type == conditionTypeRuntime &&
			condition.Status == metav1.ConditionTrue &&
			condition.Reason == reasonSyncedCoin {
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
		Reason:             reasonSyncedCoin,
		Message:            "synced coin price",
	}

	if !found {
		object.Status = cryptov1beta1.CoinStatus{
			Meta: cryptov1beta1.CoinStatusMeta{
				Price:   price.Data.Amount,
				Balance: fmt.Sprintf("%f", n*p),
			},
			Phase:      phaseRunning,
			Conditions: append(object.Status.Conditions, condition),
			Message:    "synced coin price",
			Reason:     reasonSyncedCoin,
		}
		updated = true
	} else {
		if foundIndex >= 0 {
			object.Status.Conditions[foundIndex] = condition
			object.Status = cryptov1beta1.CoinStatus{
				Meta: cryptov1beta1.CoinStatusMeta{
					Price:   price.Data.Amount,
					Balance: fmt.Sprintf("%f", n*p),
				},
				Phase:      phaseRunning,
				Conditions: object.Status.Conditions,
				Message:    "synced coin price",
				Reason:     reasonSyncedCoin,
			}
		}
		updated = true
	}

	if updated {
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

// coinprice is the response from coinbase coin price query
type coinprice struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}
