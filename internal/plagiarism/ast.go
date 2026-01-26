package plagiarism

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/RishiKendai/aegis/internal/models"
)

// ASTSimilarity calculates similarity using AST Merkle hashing
// Uses post-order traversal to build Merkle tree hashes for all subtrees
func ASTSimilarity(artifactA, artifactB *models.Artifact) float64 {
	if artifactA.AST == nil || artifactB.AST == nil {
		return 0.0
	}

	// Build multiset of subtree hashes for both ASTs
	subtreesA := buildSubtreeHashes(artifactA.AST)
	subtreesB := buildSubtreeHashes(artifactB.AST)

	// Count common subtrees
	commonCount := 0
	for hash := range subtreesA {
		if subtreesB[hash] {
			commonCount++
		}
	}

	totalA := len(subtreesA)
	totalB := len(subtreesB)

	if totalA == 0 || totalB == 0 {
		return 0.0
	}

	// ASTScore = common_subtrees / min(total_subtrees_A, total_subtrees_B)
	minTotal := totalA
	if totalB < minTotal {
		minTotal = totalB
	}

	if minTotal == 0 {
		return 0.0
	}

	return float64(commonCount) / float64(minTotal)
}

// buildSubtreeHashes builds a multiset of subtree hashes using post-order traversal
// Returns a set of all subtree hashes (Merkle tree hashes for each node and its descendants)
func buildSubtreeHashes(node *models.ASTNode) map[string]bool {
	if node == nil {
		return make(map[string]bool)
	}

	// Hash cache: maps node pointer to its computed hash
	// This avoids recomputing hashes during traversal (production-ready optimization)
	hashCache := make(map[*models.ASTNode]string)
	
	// Set to collect all subtree hashes
	subtreeHashes := make(map[string]bool)
	
	// Build hashes using post-order traversal
	buildSubtreeHashesRecursive(node, hashCache, subtreeHashes)
	
	return subtreeHashes
}

// buildSubtreeHashesRecursive recursively builds subtree hashes using post-order traversal
// Post-order ensures children are hashed before their parent, enabling Merkle tree structure
// hashCache stores computed hashes to avoid redundant recomputation
func buildSubtreeHashesRecursive(
	node *models.ASTNode,
	hashCache map[*models.ASTNode]string,
	subtreeHashes map[string]bool,
) {
	if node == nil {
		return
	}

	// Process children first (post-order traversal)
	childHashes := make([]string, 0)
	if node.Children != nil && len(node.Children) > 0 {
		for _, child := range node.Children {
			// Recursively process child first
			buildSubtreeHashesRecursive(child, hashCache, subtreeHashes)
			
			// Get cached hash of child (already computed in recursive call above)
			if childHash, exists := hashCache[child]; exists {
				childHashes = append(childHashes, childHash)
			}
		}
	}

	// Compute hash for this node using all its properties
	nodeHash := computeNodeHash(node, childHashes)
	
	// Cache the hash for this node to avoid recomputation
	hashCache[node] = nodeHash
	
	// Add this subtree hash to the collection
	subtreeHashes[nodeHash] = true
}

// computeNodeHash computes Merkle hash for a node
// childHashes should already be sorted and computed from post-order traversal
// Includes all relevant node properties for accurate similarity detection
func computeNodeHash(node *models.ASTNode, childHashes []string) string {
	if node == nil {
		return ""
	}

	var parts []string
	
	// Node type is the primary identifier
	parts = append(parts, "type:", node.Type)
	
	// Include node name if present (for identifiers, function names, etc.)
	if node.Name != "" {
		parts = append(parts, "name:", node.Name)
	}
	
	// Include return type if present (for functions, methods)
	if node.ReturnType != "" {
		parts = append(parts, "return:", node.ReturnType)
	}
	
	// Include operator if present (for binary/unary expressions)
	if node.Operator != "" {
		parts = append(parts, "op:", node.Operator)
	}
	
	// Include modifiers if present (for classes, methods, etc.)
	if len(node.Modifiers) > 0 {
		modifiers := make([]string, len(node.Modifiers))
		copy(modifiers, node.Modifiers)
		sort.Strings(modifiers) // Sort for consistency
		parts = append(parts, "modifiers:", strings.Join(modifiers, ","))
	}
	
	// Include parameters if present (for function/method declarations)
	if len(node.Parameters) > 0 {
		paramStrings := make([]string, 0, len(node.Parameters))
		for _, param := range node.Parameters {
			paramStr := fmt.Sprintf("%s:%s", param.Type, param.Name)
			if param.ParamType != "" {
				paramStr += ":" + param.ParamType
			}
			paramStrings = append(paramStrings, paramStr)
		}
		sort.Strings(paramStrings) // Sort for consistency
		parts = append(parts, "params:", strings.Join(paramStrings, ","))
	}
	
	// Include child hashes (already sorted from caller)
	// This creates the Merkle tree structure: parent hash depends on children hashes
	if len(childHashes) > 0 {
		parts = append(parts, "children:", strings.Join(childHashes, ","))
	}
	
	// Combine all parts into a single string
	hashInput := strings.Join(parts, "|")
	
	// Compute SHA256 hash
	return computeHash(hashInput)
}

// computeHash computes SHA256 hash of a string and returns hex-encoded result
func computeHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
