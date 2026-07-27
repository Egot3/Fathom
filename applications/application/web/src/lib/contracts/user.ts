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
