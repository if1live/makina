package hitomi

import "testing"
import "github.com/stretchr/testify/assert"

func TestExtractCoverImages(t *testing.T) {
	jsontext := `
    
{
  "metadata": {
    "id": "993344",
    "title": "I My Mask Ch. 1-4",
    "covers": [
      "https://tn.hitomi.la/bigtn/993344/000.jpg.jpg",
      "https://tn.hitomi.la/bigtn/993344/001.jpg.jpg"
    ],
    "artists": [
      "nanamiya tsugumi"
    ],...
`
	expected := []string{
		"https://tn.hitomi.la/bigtn/993344/000.jpg.jpg",
		"https://tn.hitomi.la/bigtn/993344/001.jpg.jpg",
	}
	actual := ExtractCoverImages(jsontext)
	assert.Equal(t, expected, actual)
}
