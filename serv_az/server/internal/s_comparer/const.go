package comparer

import "server/internal/domain/models"

// i think this is a bad way to save some data, it's better to use reddis or some hashDB
var HashTableWords = make(map[string]*models.TestPageData)
var HashTableWordsLearn = make(map[string]*models.TestPageData)
