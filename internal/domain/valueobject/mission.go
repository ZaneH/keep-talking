package valueobject

import "time"

type Mission int32

const (
	MissionUnspecified Mission = iota

	// Section 1: Introduction
	MissionTheFirstBomb

	// Section 2: The Basics
	MissionSomethingOldSomethingNew
	MissionDoubleYourMoney
	MissionOneStepUp
	MissionPickUpThePace

	// Section 3: Moderate
	MissionAHiddenMessage
	MissionSomethingsDifferent
	MissionOneGiantLeap
	MissionFairGame
	MissionPickUpThePaceII
	MissionNoRoomForError
	MissionEightMinutes

	// Section 4: Needy Modules
	MissionASmallWrinkle
	MissionPayAttention
	MissionTheKnob
	MissionMultiTasker

	// Section 5: Challenging
	MissionWiresWiresEverywhere
	MissionComputerHacking
	MissionWhosOnFirstChallenge
	MissionFiendish
	MissionPickUpThePaceIII
	MissionOneWithEverything

	// Section 6: Extreme
	MissionPickUpThePaceIV
	MissionJuggler
	MissionDoubleTrouble
	MissionIAmHardcore

	// Section 7: Exotic
	MissionBlinkenlights
	MissionAppliedTheory
	MissionAMazeIng
	MissionSnipSnap
	MissionRainbowTable
	MissionBlinkenlightsII
)

type ModuleSpec struct {
	Type          ModuleType
	PossibleTypes []ModuleType
	Count         int
}

type MissionDefinition struct {
	Name       string
	Timer      time.Duration
	MaxStrikes int
	Modules    []ModuleSpec
	NumFaces   int
	Rows       int
	Columns    int
	Section    int
}

// Section-based module pools - modules available at each section for random selection
var SectionModulePools = map[int][]ModuleType{
	1: {
		BigButtonModule,
		KeypadModule,
		WiresModule,
	},
	2: {
		BigButtonModule,
		KeypadModule,
		WiresModule,
	},
	3: {
		BigButtonModule,
		KeypadModule,
		WiresModule,
		PasswordModule,
		MorseModule,
		ComplicatedWiresModule,
		WireSequenceModule,
		SimonModule,
	},
	4: {
		BigButtonModule,
		KeypadModule,
		WiresModule,
		PasswordModule,
		MorseModule,
		ComplicatedWiresModule,
		WireSequenceModule,
		SimonModule,
		NeedyVentGasModule,
		NeedyKnobModule,
	},
	5: {
		// All modules
		BigButtonModule,
		KeypadModule,
		WiresModule,
		PasswordModule,
		MorseModule,
		ComplicatedWiresModule,
		WireSequenceModule,
		SimonModule,
		WhosOnFirstModule,
		MemoryModule,
		MazeModule,
		NeedyVentGasModule,
		NeedyKnobModule,
	},
	6: {
		// All modules (same as section 5)
		BigButtonModule,
		KeypadModule,
		WiresModule,
		PasswordModule,
		MorseModule,
		ComplicatedWiresModule,
		WireSequenceModule,
		SimonModule,
		WhosOnFirstModule,
		MemoryModule,
		MazeModule,
		NeedyVentGasModule,
		NeedyKnobModule,
	},
	7: {
		// All modules (same as sections 5 & 6)
		BigButtonModule,
		KeypadModule,
		WiresModule,
		PasswordModule,
		MorseModule,
		ComplicatedWiresModule,
		WireSequenceModule,
		SimonModule,
		WhosOnFirstModule,
		MemoryModule,
		MazeModule,
		NeedyVentGasModule,
		NeedyKnobModule,
	},
}

// Special marker for random module selection
const RandomModule ModuleType = -1

var MissionDefinitions = map[Mission]MissionDefinition{
	MissionTheFirstBomb: {
		Name:       "The First Bomb",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    1,
		Modules: []ModuleSpec{
			{Type: BigButtonModule, Count: 1},
			{Type: KeypadModule, Count: 1},
			{Type: WiresModule, Count: 1},
		},
	},
	MissionSomethingOldSomethingNew: {
		Name:       "Something Old, Something New",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    2,
		Modules: []ModuleSpec{
			{Type: KeypadModule, Count: 1},
			{Type: WiresModule, Count: 1},
			{Type: RandomModule, Count: 1},
		},
	},
	MissionDoubleYourMoney: {
		Name:       "Double Your Money",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    2,
		Modules: []ModuleSpec{
			{Type: BigButtonModule, Count: 2},
			{Type: KeypadModule, Count: 2},
			{Type: WiresModule, Count: 2},
		},
	},
	MissionOneStepUp: {
		Name:       "One Step Up",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    2,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 4},
		},
	},
	MissionPickUpThePace: {
		Name:       "Pick up the Pace",
		Timer:      3 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    2,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 3},
		},
	},
	MissionAHiddenMessage: {
		Name:       "A Hidden Message",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{PossibleTypes: []ModuleType{MorseModule, PasswordModule}, Count: 1},
			{Type: RandomModule, Count: 2},
		},
	},
	MissionSomethingsDifferent: {
		Name:       "Something's Different",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{PossibleTypes: []ModuleType{ComplicatedWiresModule, WireSequenceModule}, Count: 1},
			{Type: RandomModule, Count: 2},
		},
	},
	MissionOneGiantLeap: {
		Name:       "One Giant Leap",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 4},
		},
	},
	MissionFairGame: {
		Name:       "Fair Game",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 5},
		},
	},
	MissionPickUpThePaceII: {
		Name:       "Pick up the Pace II",
		Timer:      2*time.Minute + 30*time.Second,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 5},
		},
	},
	MissionNoRoomForError: {
		Name:       "No Room for Error",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 5},
		},
	},
	MissionEightMinutes: {
		Name:       "Eight Minutes",
		Timer:      8 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    3,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 8},
		},
	},
	MissionASmallWrinkle: {
		Name:       "A Small Wrinkle",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    4,
		Modules: []ModuleSpec{
			{Type: NeedyVentGasModule, Count: 1},
			{Type: RandomModule, Count: 5},
		},
	},
	MissionPayAttention: {
		Name:       "Pay Attention",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    4,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 4},
			{Type: NeedyVentGasModule, Count: 1},
		},
	},
	MissionTheKnob: {
		Name:       "The Knob",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    4,
		Modules: []ModuleSpec{
			{Type: NeedyKnobModule, Count: 1},
			{Type: RandomModule, Count: 5},
		},
	},
	MissionMultiTasker: {
		Name:       "Multi-tasker",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    4,
		Modules: []ModuleSpec{
			{Type: NeedyVentGasModule, Count: 1},
			{Type: NeedyKnobModule, Count: 1},
			{Type: RandomModule, Count: 4},
		},
	},
	MissionWiresWiresEverywhere: {
		Name:       "Wires! Wires Everywhere!",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: WiresModule, Count: 2},
			{Type: ComplicatedWiresModule, Count: 2},
			{Type: WireSequenceModule, Count: 2},
		},
	},
	MissionComputerHacking: {
		Name:       "Computer Hacking",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: NeedyVentGasModule, Count: 5},
			{Type: PasswordModule, Count: 1},
			{Type: SimonModule, Count: 1},
			{Type: MazeModule, Count: 1},
		},
	},
	MissionWhosOnFirstChallenge: {
		Name:       "Who's on First?",
		Timer:      3*time.Minute + 30*time.Second,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: WhosOnFirstModule, Count: 4},
		},
	},
	MissionFiendish: {
		Name:       "Fiendish",
		Timer:      5 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: NeedyVentGasModule, Count: 1},
			{Type: RandomModule, Count: 5},
		},
	},
	MissionPickUpThePaceIII: {
		Name:       "Pick up the Pace III",
		Timer:      90 * time.Second,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: RandomModule, Count: 4},
		},
	},
	MissionOneWithEverything: {
		Name:       "One with Everything",
		Timer:      6 * time.Minute,
		MaxStrikes: 3,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    5,
		Modules: []ModuleSpec{
			{Type: WiresModule, Count: 1},
			{Type: BigButtonModule, Count: 1},
			{Type: KeypadModule, Count: 1},
			{Type: SimonModule, Count: 1},
			{Type: WhosOnFirstModule, Count: 1},
			{Type: MemoryModule, Count: 1},
			{Type: MorseModule, Count: 1},
			{Type: ComplicatedWiresModule, Count: 1},
			{Type: WireSequenceModule, Count: 1},
			{Type: MazeModule, Count: 1},
			{Type: PasswordModule, Count: 1},
		},
	},
	MissionPickUpThePaceIV: {
		Name:       "Pick up the Pace IV",
		Timer:      80 * time.Second,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    6,
		Modules: []ModuleSpec{
			{Type: KeypadModule, Count: 1},
			{Type: ComplicatedWiresModule, Count: 1},
			{Type: WireSequenceModule, Count: 1},
			{PossibleTypes: []ModuleType{BigButtonModule, WiresModule}, Count: 1},
		},
	},
	MissionJuggler: {
		Name:       "Juggler",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    6,
		Modules: []ModuleSpec{
			{Type: NeedyKnobModule, Count: 1},
			{Type: SimonModule, Count: 1},
			{Type: WiresModule, Count: 1},
			{Type: MorseModule, Count: 1},
			{Type: NeedyVentGasModule, Count: 1},
			{Type: RandomModule, Count: 3},
		},
	},
	MissionDoubleTrouble: {
		Name:       "Double Trouble",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    6,
		Modules: []ModuleSpec{
			{Type: NeedyKnobModule, Count: 2},
			{Type: RandomModule, Count: 6},
		},
	},
	MissionIAmHardcore: {
		Name:       "I am Hardcore",
		Timer:      5 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    6,
		Modules: []ModuleSpec{
			{PossibleTypes: []ModuleType{NeedyKnobModule, NeedyVentGasModule}, Count: 1},
			{Type: RandomModule, Count: 10},
		},
	},
	MissionBlinkenlights: {
		Name:       "Blinkenlights",
		Timer:      90 * time.Second,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: SimonModule, Count: 5},
		},
	},
	MissionAppliedTheory: {
		Name:       "Applied Theory",
		Timer:      3 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: ComplicatedWiresModule, Count: 11},
		},
	},
	MissionAMazeIng: {
		Name:       "A-Maze-Ing",
		Timer:      3 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: MazeModule, Count: 8},
		},
	},
	MissionSnipSnap: {
		Name:       "Snip Snap",
		Timer:      3 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: WireSequenceModule, Count: 6},
		},
	},
	MissionRainbowTable: {
		Name:       "Rainbow Table",
		Timer:      4 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   2,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: PasswordModule, Count: 9},
		},
	},
	MissionBlinkenlightsII: {
		Name:       "Blinkenlights II",
		Timer:      3 * time.Minute,
		MaxStrikes: 1,
		NumFaces:   1,
		Rows:       2,
		Columns:    3,
		Section:    7,
		Modules: []ModuleSpec{
			{Type: SimonModule, Count: 3},
			{Type: MorseModule, Count: 3},
		},
	},
}
