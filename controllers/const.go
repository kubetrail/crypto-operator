package controllers

const (
	finalizer                     = "crypto.kubetrail.io/finalizer"
	label                         = "crypto.kubetrail.io/group"
	reasonObjectInitialized       = "objectInitialized"
	reasonObjectMarkedForDeletion = "objectMarkedForDeletion"
	reasonFinalizerAdded          = "finalizerAdded"
	reasonSynced                  = "synced"
	reasonDeleted                 = "deleted"
	phasePending                  = "pending"
	phaseRunning                  = "running"
	phaseError                    = "error"
	phaseTerminating              = "terminating"
	conditionTypeObject           = "object"
	conditionTypeRuntime          = "runtime"
)
