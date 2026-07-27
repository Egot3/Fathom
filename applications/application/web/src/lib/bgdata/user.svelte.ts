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
	console.log("Token expiration: ", GetTokenExpiration());
	if (GetTokenExpiration() < new Date()) {
		// window.localStorage.removeItem(UserKey); //pretty unnecessary, as this info is cosmetic and not real auth
		return null;
	}
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

function GetUserOrRedirect(): userData {
	const u = GetUser();
	if (u == null) {
		window.location.href = "/login";
		return null as never; // holy hell
	}
	return u;
}

let tokenExpiration: Date = $state(new Date(Date.now() - 1));
const TokenExpKey = "tokenExpiration";

function GetTokenExpiration(): Date {
	if (tokenExpiration < new Date()) {
		const cached = window.localStorage.getItem(TokenExpKey);
		if (cached == null) {
			return tokenExpiration;
		}

		return new Date(cached);
	}

	return tokenExpiration;
}

function SetTokenExpiration(newExp: Date) {
	console.log("new token exp: ", newExp);
	window.localStorage.setItem(TokenExpKey, newExp.toISOString());
	tokenExpiration = newExp;
}

export {
	SetUser,
	GetUser,
	GetTokenExpiration,
	SetTokenExpiration,
	GetUserOrRedirect,
};
