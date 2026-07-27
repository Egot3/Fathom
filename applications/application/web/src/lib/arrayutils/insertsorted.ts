export function InsertSorted<T>(arr: T[], value: T) {
	let low = 0;
	let high = arr.length;

	while (low < high) {
		let mid = (low + high) >> 1;
		if (arr[mid] < value) {
			low = mid + 1;
		} else {
			high = mid;
		}
	}

	arr.splice(low, 0, value);
	return arr;
}
