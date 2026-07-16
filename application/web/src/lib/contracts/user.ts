type User = {
	UUID: string;
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

export type { User, LoginRequest, LoginResponse };
