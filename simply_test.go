package simply

import "testing"

func TestSimply_Status_Equals_Success(context *testing.T) {
	// Test-ception
	testSubject := Test(context, "Subject")

	testSimply := Test(context, "Test Pending Target")
	testSimply.Target(testSubject.GetStatus())
	testSimply.Equals(PendingTarget)
	testSimply.Validate(testSimply.result)

	testSubject.Target(testSimply.result)

	testSimply = Test(context, "Test Pending Comparison")
	testSimply.Target(testSubject.GetStatus())
	testSimply.Equals(PendingComparison)
	testSimply.Validate(testSimply.result)

	testSubject.Equals(testSimply.result)

	testSimply = Test(context, "Test Pass Pending Validation")
	testSimply.Target(testSubject.GetStatus())
	testSimply.Equals(PassPendingValidation)
	testSimply.Validate(testSimply.result)

	testSubject.Validate(testSubject.result)

	testSimply = Test(context, "Test Pass Validated")
	testSimply.Target(testSubject.GetStatus())
	testSimply.Equals(Passed)
	testSimply.Validate(testSimply.result)
}
