function decrypt(a) {
    return a - Math.log10(a * Math.LN10)
}

console.log(decrypt(0.42850214618580407)) // 0.43433345602451556