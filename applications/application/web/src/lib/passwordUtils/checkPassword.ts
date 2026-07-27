export const CheckPassword = (passwordString: string): string => {
	const passwordLength = passwordString.length;
	if (passwordLength <= 8) {
		return "password must be longer than 7 characters";
	}
	if (passwordLength >= 255) {
		return "password must be shorter than 256 characters";
	}

	const digitCount = (passwordString.match(/\d/g) || []).length;
	if (digitCount < 4) {
		return "5 digits in password are required"; // in bank account too;
	}

	const whiteSpace = /\s+/.exec(passwordString);
	if (whiteSpace !== null) {
		return `Password mustn't have any whitespace characters`;
	}

	const uppercaseCount = (passwordString.match(/[A-Z]/g) || []).length;
	if (uppercaseCount < 1) {
		return "Password must have at least 2 UPPERCASE letters";
	}

	const lowercaseCount = (passwordString.match(/([a-z])/g) || []).length;
	if (lowercaseCount < 1) {
		return "Password must have at least 2 lowercase letters";
	}

	return "";
};
