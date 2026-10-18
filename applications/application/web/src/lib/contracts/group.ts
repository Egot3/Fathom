import { err, errAsync, okAsync, ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";
import {
  NormalizeJSON,
  NormalizeResponse,
  TokenizedFetch,
} from "./tokenizedFetch";
import type { User } from "./user";

export type Group = {
  uuid: string;
  name: string;

  pupils: User[];
};

export type ListGroupResponse = {
  page: number;
  size: number;
  total: number;
  groups: Group[];
};

export function FetchGroups(
  page: number,
  size: number,
): ResultAsync<ListGroupResponse, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/group?page=${page}&size=${size}`,
    {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    },
  ).andThen(NormalizeJSON<ListGroupResponse>);
}

type PostGroupRequest = {
  name: string;
  appendants?: string[];
};

export function FetchGroupPost(
  name: string,
  appendants?: string[],
): ResultAsync<null, JSONError> {
  const body: PostGroupRequest = {
    name: name,
    appendants: appendants,
  };
  const bodyString = JSON.stringify(body);
  return TokenizedFetch(
    "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/group/",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}

export function FetchGroupDelete(UUID: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/group/${UUID}`,
    {
      method: "DELETE",
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}

type PatchGroupRequest = {
  name: string | undefined;
};

export function FetchGroupPatch(
  UUID: string,
  name: string | undefined,
): ResultAsync<null, JSONError> {
  const body: PatchGroupRequest = {
    name: name,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/group/" + UUID,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}

type PruneGroupRequest = {
  removants: string[];
};

export function FetchGroupPrune(
  groupUUID: string,
  userUUIDs: string[],
): ResultAsync<null, JSONError> {
  const body: PruneGroupRequest = {
    removants: userUUIDs,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/group/${groupUUID}/user`,
    {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}

type AppendGroupRequest = {
  appendants: string[];
};

export function FetchGroupAppend(
  groupUUID: string,
  userUUIDs: string[],
): ResultAsync<null, JSONError> {
  const body: AppendGroupRequest = {
    appendants: userUUIDs,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    `https://${import.meta.env.VITE_DOMAIN}/api/v1/group/${groupUUID}/user`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: bodyString,
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}
