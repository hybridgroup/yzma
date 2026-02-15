package llama

import (
	"testing"
)

func TestLoraInit(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	if adapter == 0 {
		t.Fatalf("AdapterLoraInit returned null adapter")
	}
	defer AdapterLoraFree(adapter)

	t.Logf("LoRA adapter initialized successfully: %v", adapter)
}

func TestAdapterMetaCount(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	count := AdapterMetaCount(adapter)
	t.Logf("AdapterMetaCount returned: %d", count)
	if count < 0 {
		t.Fatal("AdapterMetaCount returned negative value")
	}
}

func TestAdapterMetaKeyByIndex(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	count := AdapterMetaCount(adapter)
	if count <= 0 {
		t.Skip("No adapter metadata keys to test")
	}
	key, ok := AdapterMetaKeyByIndex(adapter, 0)
	if !ok {
		t.Fatal("AdapterMetaKeyByIndex failed for index 0")
	}
	t.Logf("AdapterMetaKeyByIndex returned: %s", key)
}

func TestAdapterMetaValStrByIndex(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}

	count := AdapterMetaCount(adapter)
	if count <= 0 {
		t.Skip("No adapter metadata values to test")
	}
	defer AdapterLoraFree(adapter)

	val, ok := AdapterMetaValStrByIndex(adapter, 0)
	if !ok {
		t.Fatal("AdapterMetaValStrByIndex failed for index 0")
	}
	t.Logf("AdapterMetaValStrByIndex returned: %s", val)
}

func TestAdapterMetaValStr(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	count := AdapterMetaCount(adapter)
	if count <= 0 {
		t.Skip("No adapter metadata to test")
	}
	key, ok := AdapterMetaKeyByIndex(adapter, 0)
	if !ok {
		t.Skip("AdapterMetaKeyByIndex failed for index 0")
	}
	val, ok := AdapterMetaValStr(adapter, key)
	if !ok {
		t.Fatal("AdapterMetaValStr failed for key:", key)
	}
	t.Logf("AdapterMetaValStr returned: %s", val)
}

func TestSetAdaptersLora(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	ctx, err := InitFromModel(model, ContextDefaultParams())
	if err != nil {
		t.Fatalf("InitFromModel failed: %v", err)
	}
	defer Free(ctx)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	adapter2, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter2)

	ret := SetAdaptersLora(ctx, []AdapterLora{adapter, adapter2}, []float32{1.0, 1.0})
	if ret != 0 {
		t.Fatalf("SetAdaptersLora failed, return code: %d", ret)
	}
	t.Logf("SetAdaptersLora succeeded")
}

func TestAdapterGetAloraNInvocationTokens(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	n := AdapterGetAloraNInvocationTokens(adapter)
	t.Logf("AdapterGetAloraNInvocationTokens returned: %d", n)
}

func TestAdapterGetAloraInvocationTokens(t *testing.T) {
	modelFile := testLoraModelFileName(t)
	loraFile := testLoraAdaptorFileName(t)

	testSetup(t)
	defer testCleanup(t)

	model, err := ModelLoadFromFile(modelFile, ModelDefaultParams())
	if err != nil {
		t.Fatalf("ModelLoadFromFile failed: %v", err)
	}
	defer ModelFree(model)

	adapter, err := AdapterLoraInit(model, loraFile)
	if err != nil {
		t.Fatalf("AdapterLoraInit failed: %v", err)
	}
	defer AdapterLoraFree(adapter)

	tokens := AdapterGetAloraInvocationTokens(adapter)
	t.Logf("AdapterGetAloraInvocationTokens returned %d tokens", len(tokens))
}
