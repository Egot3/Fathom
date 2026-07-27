export function FindSorted<T>(arr: T[], value: T): number {
	let low = 0;
	let high = arr.length;

	while (low < high) {
		let mid = (low + high) >> 1;
		if (arr[mid] < value) {
			low = mid + 1;
		} else {
			if (arr[mid] > value) {
				high = mid;
			} else {
				return mid;
			}
		}
	}

	return -1;
}
