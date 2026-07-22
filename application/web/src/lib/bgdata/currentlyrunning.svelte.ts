import {
	FetchCurrentlyRunningQuizUUIDs,
	FetchCurrentlyRunningTestInfo,
} from "../contracts/test";
import { IsJSONError, type JSONError } from "../statuses/jsonerror";

export type ETagInfo = {
	ETag: string;
	ExpiresAt: Date;
};

type currentlyRunningData = {
	UUID: string;
	Name: string;
	QuizUUIDs: string[];
	IsPaused: boolean;
	Deadline: Date;
};

let currentlyRunning: currentlyRunningData | null = $state(null);
const CurrentlyRunningKey = "currentlyRunning";

export type NullableTestOrError = currentlyRunningData | null | JSONError;

export async function GetCurrentlyRunning(): Promise<NullableTestOrError> {
	if (currentlyRunning != null) {
		return currentlyRunning;
	}

	const testInfoResponse = await FetchCurrentlyRunningTestInfo();
	if (testInfoResponse === null) {
		currentlyRunning = null;
		return null;
	}
	if (IsJSONError(testInfoResponse)) {
		return testInfoResponse;
	}

	const cr: currentlyRunningData = {
		UUID: testInfoResponse.test.uuid,
		Name: testInfoResponse.test.name,
		Deadline: testInfoResponse.deadline,
		IsPaused: testInfoResponse.isPaused,
		QuizUUIDs: [] as string[],
	};

	const ETag = GetCurrentlyRunningCaching()?.ETag;
	const quizzesResponse = await FetchCurrentlyRunningQuizUUIDs(ETag);
	if (quizzesResponse === null) {
		currentlyRunning = null;
		return null;
	}
	if (IsJSONError(quizzesResponse)) {
		return quizzesResponse;
	}

	SetCurrentlyRunningCaching(quizzesResponse.Caching);

	cr.QuizUUIDs = quizzesResponse.UUIDs;
	currentlyRunning = cr;
	return currentlyRunning;
}

let currentlyRunningCaching: ETagInfo = $state({
	ETag: "",
	ExpiresAt: new Date(Date.now() - 1),
});
const CurrentlyRunningCaching = "currentlyRunningCaching";

function GetCurrentlyRunningCaching(): ETagInfo | null {
	if (
		currentlyRunningCaching.ExpiresAt < new Date() ||
		currentlyRunningCaching.ETag == ""
	) {
		const cached = window.localStorage.getItem(CurrentlyRunningCaching);
		if (cached == null) {
			return currentlyRunningCaching;
		}

		return JSON.parse(cached);
	}

	return currentlyRunningCaching;
}

function SetCurrentlyRunningCaching(newETagInfo: ETagInfo) {
	window.localStorage.setItem(
		CurrentlyRunningCaching,
		JSON.stringify(newETagInfo),
	);
	currentlyRunningCaching = newETagInfo;
}
