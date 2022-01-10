package controllers

const (
	finalizer                     = "crypto.kubetrail.io/finalizer"
	reasonObjectInitialized       = "objectInitialized"
	reasonObjectMarkedForDeletion = "objectMarkedForDeletion"
	reasonFinalizerAdded          = "finalizerAdded"
	reasonSyncedCoin              = "syncedCoin"
	reasonDeletedCoin             = "deletedCoin"
	phasePending                  = "pending"
	phaseRunning                  = "running"
	phaseTerminating              = "terminating"
	conditionTypeObject           = "object"
	conditionTypeRuntime          = "runtime"
)
