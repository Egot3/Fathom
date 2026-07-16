export enum InputStatus {
	Idle = "idle",
	Punish = "punish",
	Treat = "treat",
}

export function ClassForStatus(status: InputStatus): string {
	switch (status) {
		case InputStatus.Idle:
			return "";
		case InputStatus.Punish:
			return "border-error-300";
		case InputStatus.Treat:
			return "border-success-300";
	}
}
