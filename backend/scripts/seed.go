package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joaosantos/jlpt5/internal/config"
	"github.com/joaosantos/jlpt5/internal/infrastructure/database"
	"github.com/joaosantos/jlpt5/internal/utils"
)

func main() {
	fmt.Println("🌱 Starting database seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create logger
	logger := utils.NewLogger("INFO")

	// Connect to database
	db, err := database.NewPostgresConnection(&cfg.Database, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Seed vocabulary
	fmt.Println("\n📚 Seeding vocabulary...")
	if err := seedVocabulary(ctx, db); err != nil {
		log.Fatalf("Failed to seed vocabulary: %v", err)
	}
	fmt.Println("✅ Vocabulary seeded successfully")

	// Seed grammar
	fmt.Println("\n📖 Seeding grammar lessons...")
	if err := seedGrammar(ctx, db); err != nil {
		log.Fatalf("Failed to seed grammar: %v", err)
	}
	fmt.Println("✅ Grammar seeded successfully")

	// Seed quizzes
	fmt.Println("\n📝 Seeding quizzes...")
	if err := seedQuizzes(ctx, db); err != nil {
		log.Fatalf("Failed to seed quizzes: %v", err)
	}
	fmt.Println("✅ Quizzes seeded successfully")

	fmt.Println("\n🎉 Database seeding completed successfully!")
	os.Exit(0)
}

func seedVocabulary(ctx context.Context, db *database.DB) error {
	vocabulary := []struct {
		Word        string
		Reading     string
		Meaning     string
		PartOfSpeech string
		Example     *string
	}{
		{"私", "わたし", "I, me", "pronoun", strPtr("私は学生です。(I am a student.)")},
		{"あなた", "あなた", "you", "pronoun", strPtr("あなたは先生ですか。(Are you a teacher?)")},
		{"これ", "これ", "this", "pronoun", strPtr("これは本です。(This is a book.)")},
		{"それ", "それ", "that", "pronoun", strPtr("それはペンです。(That is a pen.)")},
		{"ここ", "ここ", "here", "noun", strPtr("ここは学校です。(This is a school.)")},
		{"そこ", "そこ", "there", "noun", strPtr("そこは図書館です。(That is a library.)")},
		{"今", "いま", "now", "noun", strPtr("今は三時です。(It's 3 o'clock now.)")},
		{"昨日", "きのう", "yesterday", "noun", strPtr("昨日は雨でした。(It was rainy yesterday.)")},
		{"今日", "きょう", "today", "noun", strPtr("今日は晴れです。(It's sunny today.)")},
		{"明日", "あした", "tomorrow", "noun", strPtr("明日は月曜日です。(Tomorrow is Monday.)")},
		{"学校", "がっこう", "school", "noun", strPtr("学校に行きます。(I go to school.)")},
		{"先生", "せんせい", "teacher", "noun", strPtr("田中先生は親切です。(Teacher Tanaka is kind.)")},
		{"学生", "がくせい", "student", "noun", strPtr("私は学生です。(I am a student.)")},
		{"友達", "ともだち", "friend", "noun", strPtr("友達と遊びます。(I play with friends.)")},
		{"本", "ほん", "book", "noun", strPtr("本を読みます。(I read books.)")},
		{"食べる", "たべる", "to eat", "verb", strPtr("朝ごはんを食べます。(I eat breakfast.)")},
		{"飲む", "のむ", "to drink", "verb", strPtr("水を飲みます。(I drink water.)")},
		{"見る", "みる", "to see, to watch", "verb", strPtr("テレビを見ます。(I watch TV.)")},
		{"行く", "いく", "to go", "verb", strPtr("学校に行きます。(I go to school.)")},
		{"来る", "くる", "to come", "verb", strPtr("友達が来ます。(A friend is coming.)")},
	}

	for _, v := range vocabulary {
		query := `
			INSERT INTO vocabulary (word, reading, meaning, part_of_speech, example_sentence, jlpt_level)
			VALUES ($1, $2, $3, $4, $5, 5)
			ON CONFLICT DO NOTHING
		`
		_, err := db.ExecContext(ctx, query, v.Word, v.Reading, v.Meaning, v.PartOfSpeech, v.Example)
		if err != nil {
			return fmt.Errorf("error inserting vocabulary '%s': %w", v.Word, err)
		}
	}

	return nil
}

func seedGrammar(ctx context.Context, db *database.DB) error {
	// Insert grammar lessons
	lessons := []struct {
		Title        string
		GrammarPoint string
		Explanation  string
		UsageNotes   *string
		Examples     []struct {
			Japanese string
			English  string
			Notes    *string
		}
	}{
		{
			Title:        "Basic Sentence Structure: XはYです",
			GrammarPoint: "XはYです",
			Explanation:  "This is the most basic sentence pattern in Japanese. は (wa) is the topic marker and です (desu) is the copula meaning 'is/am/are'. Use this pattern to state that X is Y.",
			UsageNotes:   strPtr("Remember that は is pronounced 'wa' when used as a particle, not 'ha'."),
			Examples: []struct {
				Japanese string
				English  string
				Notes    *string
			}{
				{"私は学生です。", "I am a student.", nil},
				{"これは本です。", "This is a book.", nil},
				{"田中さんは先生です。", "Tanaka-san is a teacher.", nil},
			},
		},
		{
			Title:        "Question Particle: か",
			GrammarPoint: "か",
			Explanation:  "Add か (ka) to the end of a sentence to make it a question. The word order stays the same as a statement.",
			UsageNotes:   strPtr("In casual speech, か can be omitted and the question is indicated by rising intonation."),
			Examples: []struct {
				Japanese string
				English  string
				Notes    *string
			}{
				{"これは本ですか。", "Is this a book?", nil},
				{"あなたは学生ですか。", "Are you a student?", nil},
				{"田中さんは先生ですか。", "Is Tanaka-san a teacher?", nil},
			},
		},
		{
			Title:        "Negative Form: じゃありません",
			GrammarPoint: "じゃありません / ではありません",
			Explanation:  "To make a negative statement, replace です with じゃありません (casual) or ではありません (formal). Both mean 'is not / am not / are not'.",
			UsageNotes:   strPtr("じゃありません is more common in everyday conversation."),
			Examples: []struct {
				Japanese string
				English  string
				Notes    *string
			}{
				{"私は学生じゃありません。", "I am not a student.", nil},
				{"これは本ではありません。", "This is not a book.", strPtr("Formal version")},
				{"田中さんは先生じゃありません。", "Tanaka-san is not a teacher.", nil},
			},
		},
		{
			Title:        "Location Particle: に",
			GrammarPoint: "に (location/time)",
			Explanation:  "The particle に (ni) marks the location where something exists or the time when something happens. It often translates to 'at', 'in', 'on', or 'to' in English.",
			UsageNotes:   strPtr("Use に with existence verbs like います and あります, and with movement verbs like 行きます."),
			Examples: []struct {
				Japanese string
				English  string
				Notes    *string
			}{
				{"学校に行きます。", "I go to school.", nil},
				{"東京に住んでいます。", "I live in Tokyo.", nil},
				{"三時に会いましょう。", "Let's meet at 3 o'clock.", nil},
			},
		},
		{
			Title:        "Object Marker: を",
			GrammarPoint: "を",
			Explanation:  "The particle を (wo/o) marks the direct object of a sentence - the thing that receives the action of the verb.",
			UsageNotes:   strPtr("を is pronounced 'o' not 'wo', even though it's written with the 'wo' character."),
			Examples: []struct {
				Japanese string
				English  string
				Notes    *string
			}{
				{"本を読みます。", "I read a book.", nil},
				{"水を飲みます。", "I drink water.", nil},
				{"テレビを見ます。", "I watch TV.", nil},
			},
		},
	}

	for lessonIdx, lesson := range lessons {
		// Insert lesson
		var lessonID int
		lessonQuery := `
			INSERT INTO grammar_lessons (title, grammar_point, explanation, usage_notes, jlpt_level, lesson_order)
			VALUES ($1, $2, $3, $4, 5, $5)
			RETURNING id
		`
		err := db.QueryRowContext(ctx, lessonQuery, lesson.Title, lesson.GrammarPoint, lesson.Explanation, lesson.UsageNotes, lessonIdx+1).Scan(&lessonID)
		if err != nil {
			return fmt.Errorf("error inserting grammar lesson '%s': %w", lesson.Title, err)
		}

		// Insert examples for this lesson
		for exampleIdx, example := range lesson.Examples {
			exampleQuery := `
				INSERT INTO grammar_examples (grammar_lesson_id, japanese_sentence, english_translation, notes, example_order)
				VALUES ($1, $2, $3, $4, $5)
			`
			_, err := db.ExecContext(ctx, exampleQuery, lessonID, example.Japanese, example.English, example.Notes, exampleIdx+1)
			if err != nil {
				return fmt.Errorf("error inserting grammar example for lesson %d: %w", lessonID, err)
			}
		}
	}

	return nil
}

func seedQuizzes(ctx context.Context, db *database.DB) error {
	quizzes := []struct {
		Title        string
		Description  string
		QuizType     string
		PassingScore int
		Questions    []struct {
			QuestionText  string
			QuestionType  string
			CorrectAnswer string
			OptionA       string
			OptionB       string
			OptionC       string
			OptionD       string
			Explanation   *string
			Points        int
		}
	}{
		{
			Title:        "Basic Vocabulary Quiz",
			Description:  "Test your knowledge of basic JLPT N5 vocabulary",
			QuizType:     "vocabulary",
			PassingScore: 70,
			Questions: []struct {
				QuestionText  string
				QuestionType  string
				CorrectAnswer string
				OptionA       string
				OptionB       string
				OptionC       string
				OptionD       string
				Explanation   *string
				Points        int
			}{
				{
					QuestionText:  "What does '私' (わたし) mean?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "A",
					OptionA:       "I, me",
					OptionB:       "You",
					OptionC:       "He, she",
					OptionD:       "We",
					Explanation:   strPtr("私 (わたし) is the most common way to say 'I' or 'me' in Japanese."),
					Points:        1,
				},
				{
					QuestionText:  "What does '学校' (がっこう) mean?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "B",
					OptionA:       "Teacher",
					OptionB:       "School",
					OptionC:       "Student",
					OptionD:       "Book",
					Explanation:   strPtr("学校 (がっこう) means 'school'."),
					Points:        1,
				},
				{
					QuestionText:  "What does '食べる' (たべる) mean?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "C",
					OptionA:       "To drink",
					OptionB:       "To see",
					OptionC:       "To eat",
					OptionD:       "To go",
					Explanation:   strPtr("食べる (たべる) is a verb meaning 'to eat'."),
					Points:        1,
				},
				{
					QuestionText:  "What does '今日' (きょう) mean?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "B",
					OptionA:       "Yesterday",
					OptionB:       "Today",
					OptionC:       "Tomorrow",
					OptionD:       "Now",
					Explanation:   strPtr("今日 (きょう) means 'today'."),
					Points:        1,
				},
				{
					QuestionText:  "What does '友達' (ともだち) mean?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "D",
					OptionA:       "Family",
					OptionB:       "Teacher",
					OptionC:       "Student",
					OptionD:       "Friend",
					Explanation:   strPtr("友達 (ともだち) means 'friend'."),
					Points:        1,
				},
			},
		},
		{
			Title:        "Basic Grammar Quiz",
			Description:  "Test your understanding of basic JLPT N5 grammar patterns",
			QuizType:     "grammar",
			PassingScore: 70,
			Questions: []struct {
				QuestionText  string
				QuestionType  string
				CorrectAnswer string
				OptionA       string
				OptionB       string
				OptionC       string
				OptionD       string
				Explanation   *string
				Points        int
			}{
				{
					QuestionText:  "Complete: 私___学生です。(I am a student.)",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "A",
					OptionA:       "は",
					OptionB:       "が",
					OptionC:       "を",
					OptionD:       "に",
					Explanation:   strPtr("は (wa) is the topic particle used in basic 'X is Y' sentences."),
					Points:        2,
				},
				{
					QuestionText:  "How do you make a question in Japanese?",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "C",
					OptionA:       "Add ね at the end",
					OptionB:       "Change the word order",
					OptionC:       "Add か at the end",
					OptionD:       "Add よ at the end",
					Explanation:   strPtr("Add か (ka) at the end of a sentence to make it a question."),
					Points:        2,
				},
				{
					QuestionText:  "Complete: 本___読みます。(I read a book.)",
					QuestionType:  "multiple_choice",
					CorrectAnswer: "B",
					OptionA:       "は",
					OptionB:       "を",
					OptionC:       "に",
					OptionD:       "で",
					Explanation:   strPtr("を marks the direct object of the verb."),
					Points:        2,
				},
			},
		},
	}

	for _, quiz := range quizzes {
		// Insert quiz
		var quizID int
		quizQuery := `
			INSERT INTO quizzes (title, description, quiz_type, jlpt_level, passing_score)
			VALUES ($1, $2, $3, 5, $4)
			RETURNING id
		`
		err := db.QueryRowContext(ctx, quizQuery, quiz.Title, quiz.Description, quiz.QuizType, quiz.PassingScore).Scan(&quizID)
		if err != nil {
			return fmt.Errorf("error inserting quiz '%s': %w", quiz.Title, err)
		}

		// Insert questions
		for questionIdx, question := range quiz.Questions {
			questionQuery := `
				INSERT INTO quiz_questions (quiz_id, question_text, question_type, correct_answer,
					option_a, option_b, option_c, option_d, explanation, points, question_order)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			`
			_, err := db.ExecContext(ctx, questionQuery,
				quizID, question.QuestionText, question.QuestionType, question.CorrectAnswer,
				question.OptionA, question.OptionB, question.OptionC, question.OptionD,
				question.Explanation, question.Points, questionIdx+1,
			)
			if err != nil {
				return fmt.Errorf("error inserting question for quiz %d: %w", quizID, err)
			}
		}
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}
