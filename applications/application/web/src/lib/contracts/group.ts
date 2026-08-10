import type { JSONError } from "../statuses/jsonerror";
import { TokenizedFetch } from "./tokenizedFetch";
import type { User } from "./user";

export type Group = {
  uuid: string;
  name: string;

  pupils: User[];
};

type ListGroupResponse = {
  page: number;
  size: number;
  total: number;
  groups: Group[];
};

export type GroupsOrJSONError = ListGroupResponse | JSONError;

export async function FetchGroups(
  page: number,
  size: number,
): Promise<GroupsOrJSONError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/group?page=${page}&size=${size}`,
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
      },
    );

    if (!rawRes.ok) {
      return (await rawRes.json()) as JSONError;
    }

    return (await rawRes.json()) as ListGroupResponse;
  } catch (err) {
    console.log("Couldn't fetch groups: ", err);
    if (err instanceof Error) {
      return { error: "couldn't fetch groups because of in-browser error" };
    }

    return { error: "couldn't fetch groups because of unknown error" };
  }
}

type PostGroupRequest = {
  name: string;
  appendants?: string[];
};

export async function FetchGroupPost(
  name: string,
  appendants?: string[],
): Promise<null | JSONError> {
  const body: PostGroupRequest = {
    name: name,
    appendants: appendants,
  };
  const bodyString = JSON.stringify(body);
  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/group/",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch group post: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send group post because of in-browser error" };
    }

    return { error: "couldn't send group post because of unknown error" };
  }
}

export async function FetchGroupDelete(
  UUID: string,
): Promise<null | JSONError> {
  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/group/" + UUID,
      {
        method: "DELETE",
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch group delete for quiz: ", err);
    if (err instanceof Error) {
      return {
        error: "couldn't send group delete because of in-browser error",
      };
    }

    return { error: "couldn't send group delete because of unknown error" };
  }
}

type PatchGroupRequest = {
  name: string | undefined;
};

export async function FetchGroupPatch(
  UUID: string,
  name: string | undefined,
): Promise<JSONError | null> {
  const body: PatchGroupRequest = {
    name: name,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/group/" + UUID,
      {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch group patch: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send group patch because of in-browser error" };
    }

    return { error: "couldn't send group patch because of unknown error" };
  }
}

type PruneGroupRequest = {
  removants: string[];
};

export async function FetchGroupPrune(
  groupUUID: string,
  userUUIDs: string[],
): Promise<null | JSONError> {
  const body: PruneGroupRequest = {
    removants: userUUIDs,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/test/${groupUUID}/user`,
      {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch group prune: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send group prune because of in-browser error" };
    }

    return { error: "couldn't send group prune because of unknown error" };
  }
}

type AppendGroupRequest = {
  removants: string[];
};

export async function FetchGroupAppend(
  groupUUID: string,
  userUUIDs: string[],
): Promise<null | JSONError> {
  const body: AppendGroupRequest = {
    removants: userUUIDs,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/group/${groupUUID}/user`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: bodyString,
      },
    );

    if (!response.ok) {
      console.log("response is not ok!");
      return (await response.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch remove user from group: ", err);
    if (err instanceof Error) {
      return { error: "couldn't remove user because of in-browser error" };
    }

    return { error: "couldn't remove user because of unknown error" };
  }
}
