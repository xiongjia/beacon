package slicex

func Filter[T any, Slice ~[]T](src Slice, predicate func(itr T, idx int) bool) Slice {
	result := make(Slice, 0, len(src))
	for i := range src {
		if predicate(src[i], i) {
			result = append(result, src[i])
		}
	}
	return result
}

func MutFilter[T any, Slice ~[]T](mutSrc Slice, predicate func(itr T) bool) Slice {
	cnt := 0
	for _, item := range mutSrc {
		if predicate(item) {
			mutSrc[cnt] = item
			cnt++
		}
	}
	return mutSrc[:cnt]
}

func Uniq[T comparable, Slice ~[]T](src Slice) Slice {
	result := make(Slice, 0, len(src))
	allData := make(map[T]struct{}, len(src))
	for i := range src {
		if _, ok := allData[src[i]]; ok {
			continue
		}
		allData[src[i]] = struct{}{}
		result = append(result, src[i])
	}
	return result
}

func UniqBy[T any, U comparable, Slice ~[]T](src Slice, iteratee func(item T) U) Slice {
	result := make(Slice, 0, len(src))
	seen := make(map[U]struct{}, len(src))
	for i := range src {
		key := iteratee(src[i])

		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, src[i])
	}
	return result
}
