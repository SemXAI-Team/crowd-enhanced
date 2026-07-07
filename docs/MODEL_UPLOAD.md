# DeepStream Model Upload Guide

This guide details how AI engineers can package and upload their models to be used by the CrowdGuard inference backend.

## Supported Formats

The inference backend leverages NVIDIA DeepStream (`nvinfer`). It accepts the following model formats:
- `.etlt` (NVIDIA TAO Toolkit encrypted models)
- `.engine` / `.plan` (TensorRT engines)
- `.onnx` (Open Neural Network Exchange)
- `.caffemodel` / `.prototxt`

## Required Files

When uploading a model via the dashboard, you MUST include the following files:

1. **Model Weights:** The actual model file (e.g., `resnet18_detector.etlt` or `model.onnx`).
2. **Labels File:** A plain text file (e.g., `labels.txt`) containing the class names, one per line.
3. **Primary Inference Config:** A text file (e.g., `config_infer_primary.txt`) containing the DeepStream `nvinfer` configuration. 
   - *Important:* Ensure paths in the config (e.g., `model-file`, `labelfile-path`) are relative or use the generic upload path. The backend will automatically map these when starting the pipeline.
4. **Calibration File (Optional):** If using INT8 precision, provide the calibration table (e.g., `cal_trt.bin`).

## Upload Process

1. Go to **Models -> Upload Model** in the dashboard.
2. Select all the required files simultaneously in the file picker.
3. Click **Upload**.

The Go backend REST API will receive these files and place them in the configured `DEEPSTREAM_MODELS_DIR`. The active DeepStream pipeline will be restarted (or dynamically updated if using DeepStream's engine update feature) to use the new model.
