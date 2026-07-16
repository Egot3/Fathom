import { createContext } from "svelte";

type User = {
	readonly UUID: string;
	readonly nickname: string;
	readonly is_teacher: boolean;
	readonly created_at: string;
	readonly updated_at: string;
};

export const [user, setUser] = createContext<User>();
