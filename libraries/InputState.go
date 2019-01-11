package libraries

type InputState struct {
	MovingXL, MovingXR, MovingYD, MovingYU, Selecting, Running bool
}

func GenerateInputState() InputState {
	inputState := InputState{
		MovingXL: false,
		MovingXR: false,
		MovingYD: false,
		MovingYU: false,
		Selecting: false,
		Running: true }

	return inputState;
}
