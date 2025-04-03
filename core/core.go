package core

// 3つの整数のうち、最小のものを返すヘルパー関数
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// 2つの文字列 s1, s2 のレーベンシュタイン距離を返す
func levenshtein(s1, s2 string) int {
	m := len(s1)
	n := len(s2)

	// どちらかの文字列が空なら、もう一方の長さが距離
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// dp は (m+1) x (n+1) の 2次元スライスを初期化
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 初期条件：空文字列との比較
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// DP テーブルの更新
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			dp[i][j] = min3(
				dp[i-1][j]+1,      // 削除
				dp[i][j-1]+1,      // 挿入
				dp[i-1][j-1]+cost, // 置換
			)
		}
	}

	return dp[m][n]
}

// レーベンシュタイン距離に基づいた類似度を 0.0～1.0 の範囲で返す
func LevenshteinSimilarity(s1, s2 string) float64 {
	distance := levenshtein(s1, s2)
	maxLen := mmax(len(s1), len(s2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(distance)/float64(maxLen)
}

// mmax は 2 つの整数のうち大きい方を返すヘルパー関数
func mmax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
