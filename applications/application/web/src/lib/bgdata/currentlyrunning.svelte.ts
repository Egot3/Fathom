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

export function SetCurrentlyRunning(newCRD: currentlyRunningData) {
  window.localStorage.setItem(CurrentlyRunningKey, JSON.stringify(newCRD));
  currentlyRunning = newCRD;
}

export function GetCurrentlyRunning(): currentlyRunningData | null {
  if (currentlyRunning === null) {
    const cached = window.localStorage.getItem(CurrentlyRunningKey);
    console.log("cached: ", cached);
    if (cached === null) {
      return currentlyRunning;
    }

    return JSON.parse(cached);
  }

  return currentlyRunning;
}

let currentlyRunningCaching: ETagInfo = $state({
  ETag: "",
  ExpiresAt: new Date(Date.now() - 1),
});
const CurrentlyRunningCaching = "currentlyRunningCaching";

export function GetCurrentlyRunningCaching(): ETagInfo | null {
  if (
    currentlyRunningCaching.ExpiresAt < new Date() ||
    currentlyRunningCaching.ETag == ""
  ) {
    const cached = window.localStorage.getItem(CurrentlyRunningCaching);
    if (cached === null) {
      return currentlyRunningCaching;
    }

    return JSON.parse(cached);
  }

  return currentlyRunningCaching;
}

export function SetCurrentlyRunningCaching(newETag: ETagInfo) {
  window.localStorage.setItem(CurrentlyRunningCaching, JSON.stringify(newETag));
  currentlyRunningCaching = newETag;
}
