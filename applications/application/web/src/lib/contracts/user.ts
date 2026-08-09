import type { Users } from "@lucide/svelte";
import type { JSONError } from "../statuses/jsonerror";
import { TokenizedFetch } from "./tokenizedFetch";

type User = {
  uuid: string;
  nickname: string;
  is_teacher: boolean;
  created_at: string;
  updated_at: string;
};

type LoginRequest = {
  password: string;
  nickname: string;
};

type LoginResponse = {
  user: User;
};

type RegisterRequest = {
  password: string;
  nickname: string;
};

type RegisterResponse = {
  user: User;
};

type LoginResponseOrJSONError = LoginResponse | JSONError;

export async function FetchLogin(
  nickname: string,
  password: string,
): Promise<LoginResponseOrJSONError> {
  const body: LoginRequest = {
    nickname: nickname,
    password: password,
  };
  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/user/login",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
        body: JSON.stringify(body),
      },
    );

    if (!response.ok) {
      return (await response.json()) as JSONError;
    }

    return (await response.json()) as LoginResponse;
  } catch (err) {
    console.log("Couldn't fetch login for user: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send login because of in-browser error" };
    }

    return { error: "couldn't send login because of unknown error" };
  }
}

export type {
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
};

type ListUsersResponse = {
  page: number;
  size: number;
  total: number;
  users: User[];
};

export type UsersOrJSONError = ListUsersResponse | JSONError;

export async function FetchUsers(
  page: number,
  size: number,
): Promise<UsersOrJSONError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/user/?page=${page}&size=${size}`,
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

    return (await rawRes.json()) as ListUsersResponse;
  } catch (err) {
    console.log("Couldn't fetch users: ", err);
    if (err instanceof Error) {
      return { error: "couldn't fetch users because of in-browser error" };
    }

    return { error: "couldn't fetch users because of unknown error" };
  }
}

export async function FetchUserDelete(UUID: string): Promise<null | JSONError> {
  try {
    const rawRes = await TokenizedFetch(
      `https://${import.meta.env.VITE_DOMAIN}/api/v1/user/${UUID}`,
      {
        method: "DELETE",
        headers: {
          Accept: "application/json",
        },
      },
    );

    if (!rawRes.ok) {
      return (await rawRes.json()) as JSONError;
    }

    return null;
  } catch (err) {
    console.log("Couldn't fetch delete user: ", err);
    if (err instanceof Error) {
      return {
        error: "couldn't fetch delete user because of in-browser error",
      };
    }

    return { error: "couldn't fetch delete user because of unknown error" };
  }
}

type PatchUserRequest = {
  name: string | undefined;
  password: string | undefined;
  is_teacher: boolean | undefined;
};

export type PatchUserStruct = {
  Name?: string;
  Password?: string;
  IsTeacher?: boolean;
};

export async function FetchUserPatch(
  UUID: string,
  { Name, Password, IsTeacher }: PatchUserStruct,
): Promise<null | JSONError> {
  const body: PatchUserRequest = {
    name: Name,
    password: Password,
    is_teacher: IsTeacher,
  };

  const bodyString = JSON.stringify(body);

  try {
    const response = await TokenizedFetch(
      "https://" + import.meta.env.VITE_DOMAIN + "/api/v1/user/" + UUID,
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
    console.log("Couldn't fetch user patch: ", err);
    if (err instanceof Error) {
      return { error: "couldn't send user patch because of in-browser error" };
    }

    return { error: "couldn't send user patch because of unknown error" };
  }
}
