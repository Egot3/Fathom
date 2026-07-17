type userData = {
	UUID: string;
	nickname: string;
	isTeacher: boolean;
};

let user: userData | null = $state(null);
const UserKey = "user";

function SetUser(newUser: userData) {
	console.log("setting user to: ", newUser);
	window.localStorage.setItem(UserKey, JSON.stringify(newUser));
	user = newUser;
}

function GetUser(): userData | null {
	console.log("getting user");
	if (user == null) {
		console.log("user is null");
		const cached = window.localStorage.getItem(UserKey);
		console.log("cached: ", cached);
		if (cached == null) {
			return user;
		}

		return JSON.parse(cached); // practically *could* throw an error, we will solve this problem by ignoring it
	}
	return user;
}

let tokenExpiration: Date = $state(new Date(Date.now() - 1));
const ExpKey = "tokenExpiration";

function GetTokenExpiration(): Date {
	if (tokenExpiration < new Date()) {
		const cached = window.localStorage.getItem(ExpKey);
		if (cached == null) {
			return tokenExpiration;
		}

		return new Date(cached);
	}

	return tokenExpiration;
}

function SetTokenExpiration(newExp: Date) {
	window.localStorage.setItem(ExpKey, newExp.toISOString());
	tokenExpiration = newExp;
}

export { SetUser, GetUser, GetTokenExpiration, SetTokenExpiration };
