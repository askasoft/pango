package ldt

const maxTrigramDistance = 300
const maxTotalDistance = 90000

// Double maxTrigramDistance
const textTrigramSize = maxTrigramDistance * 2

// ReliableConfidenceThreshold is confidence rating that has to be succeeded
// for the language detection to be considered reliable.
const ReliableConfidenceThreshold = 0.8
