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

export type {
	User,
	LoginRequest,
	LoginResponse,
	RegisterRequest,
	RegisterResponse,
};
