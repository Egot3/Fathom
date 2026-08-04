export const CopyToClipboard = async (text: string): Promise<void> => {
	try {
		await navigator.clipboard.writeText(text);
		console.log("Text successfully copied to clipboard");
	} catch (error) {
		console.error("Failed to copy text: ", error);
	}
};
