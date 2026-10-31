import { errAsync, okAsync, ResultAsync } from "neverthrow";
import type { JSONError } from "../statuses/jsonerror";
import {
  NormalizeJSON,
  NormalizeResponse,
  TokenizedFetch,
} from "./tokenizedFetch";

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

export function FetchLogin(
  nickname: string,
  password: string,
): ResultAsync<LoginResponse, JSONError> {
  const body: LoginRequest = {
    nickname: nickname,
    password: password,
  };
  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/login`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(body),
    },
    true,
  ).andThen(NormalizeJSON<LoginResponse>);
}

export function FetchRegister(
  nickname: string,
  password: string,
): ResultAsync<RegisterResponse, JSONError> {
  const body: LoginRequest = {
    nickname: nickname,
    password: password,
  };
  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/register`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(body),
    },
    true,
  ).andThen(NormalizeJSON<RegisterResponse>);
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

export function FetchUsers(
  page: number,
  size: number,
): ResultAsync<ListUsersResponse, JSONError> {
  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/?page=${page}&size=${size}`,
    {
      method: "GET",
      headers: {
        Accept: "application/json",
      },
    },
  ).andThen(NormalizeJSON<ListUsersResponse>);
}

export function FetchUserDelete(UUID: string): ResultAsync<null, JSONError> {
  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/${UUID}`,
    {
      method: "DELETE",
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
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

export function FetchUserPatch(
  UUID: string,
  { Name, Password, IsTeacher }: PatchUserStruct,
): ResultAsync<null, JSONError> {
  const body: PatchUserRequest = {
    name: Name,
    password: Password,
    is_teacher: IsTeacher,
  };

  const bodyString = JSON.stringify(body);

  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/${UUID}`,
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

export function FetchUserImport(file: File): ResultAsync<null, JSONError> {
  const formData = new FormData();

  formData.append("imported", file, file.name);

  return TokenizedFetch(
    `${import.meta.env.PUBLIC_BASE_URL}/api/v1/user/import`,
    {
      method: "POST",
      body: formData,
    },
  )
    .andThen(NormalizeResponse)
    .map((_) => null);
}
