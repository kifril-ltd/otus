package hw03frequencyanalysis

import (
	"errors"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestFrequencyAnalyser(t *testing.T) {
	fruitsText := "apple orange apple banana"

	t.Run("empty text", func(t *testing.T) {
		analyser := NewFrequencyAnalyser()
		text := ""
		expectedTop := []string{}

		analyser.Analyse(text)
		topWords := analyser.Top(2)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("error normalized func", func(t *testing.T) {
		analyser := NewFrequencyAnalyser()
		expectedTop := []string{}

		// Добавляем нормализующую функцию, которая всегда возвращает ошибку
		analyser.AddNormalizeFunc(func(s string) (string, error) {
			return "", errors.New("mocked normalization error")
		})

		analyser.Analyse(fruitsText)
		topWords := analyser.Top(2)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("simple normalise func", func(t *testing.T) {
		analyser := NewFrequencyAnalyser()
		expectedTop := []string{"APPLE", "BANANA"}

		analyser.AddNormalizeFunc(func(s string) (string, error) {
			return strings.ToUpper(s), nil
		})

		analyser.Analyse(fruitsText)
		topWords := analyser.Top(2)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("all words in top", func(t *testing.T) {
		analyser := NewFrequencyAnalyser()
		expectedTop := []string{"apple", "banana", "orange"}

		analyser.Analyse(fruitsText)
		topWords := analyser.Top(10)

		require.Equal(t, expectedTop, topWords)
	})
}

func TestFrequencyAnalyserWithText(t *testing.T) {
	text := `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

	t.Run("without normalizers", func(t *testing.T) {
		expectedTop := []string{
			"он",        // 8
			"а",         // 6
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"-",         // 4
			"Кристофер", // 4
			"если",      // 4
			"не",        // 4
			"то",        // 4
		}

		analyser := NewFrequencyAnalyser()
		analyser.Analyse(text)
		topWords := analyser.Top(10)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("with task normalizers", func(t *testing.T) {
		expectedTop := []string{
			"а",         // 8
			"он",        // 8
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"в",         // 4
			"его",       // 4
			"если",      // 4
			"кристофер", // 4
			"не",        // 4
		}

		analyser := NewFrequencyAnalyser()
		analyser.AddNormalizeFunc([]NormalizeFunc{
			func(s string) (string, error) {
				return strings.TrimFunc(s, func(r rune) bool {
					return !unicode.IsLetter(r)
				}), nil
			},
			func(s string) (string, error) {
				if len(s) == 0 {
					return "", nil
				}

				runes := []rune(s)
				runes[0] = unicode.ToLower(runes[0])
				runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])

				return string(runes), nil
			},
		}...)

		analyser.Analyse(text)
		topWords := analyser.Top(10)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("with task normalizers top n", func(t *testing.T) {
		expectedTop := []string{
			"а",         // 8
			"он",        // 8
			"и",         // 6
			"ты",        // 5
			"что",       // 5
			"в",         // 4
			"его",       // 4
			"если",      // 4
			"кристофер", // 4
			"не",        // 4
		}

		analyser := NewFrequencyAnalyser()
		analyser.AddNormalizeFunc([]NormalizeFunc{
			func(s string) (string, error) {
				return strings.TrimFunc(s, func(r rune) bool {
					return !unicode.IsLetter(r)
				}), nil
			},
			func(s string) (string, error) {
				if len(s) == 0 {
					return "", nil
				}

				runes := []rune(s)
				runes[0] = unicode.ToLower(runes[0])
				runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])

				return string(runes), nil
			},
		}...)
		analyser.Analyse(text)

		for i := 0; i < 10; i++ {
			topWords := analyser.Top(i)
			require.Equal(t, expectedTop[:i], topWords)
		}
	})
}

func TestFrequencyAnalyserWithHyphenText(t *testing.T) {
	text := "- - - - - -"

	t.Run("without normalizers", func(t *testing.T) {
		expectedTop := []string{"-"}

		analyser := NewFrequencyAnalyser()
		analyser.Analyse(text)
		topWords := analyser.Top(10)

		require.Equal(t, expectedTop, topWords)
	})

	t.Run("with task normalizers", func(t *testing.T) {
		expectedTop := []string{}

		analyser := NewFrequencyAnalyser()
		analyser.AddNormalizeFunc([]NormalizeFunc{
			func(s string) (string, error) {
				return strings.TrimFunc(s, func(r rune) bool {
					return !unicode.IsLetter(r)
				}), nil
			},
			func(s string) (string, error) {
				if len(s) == 0 {
					return "", nil
				}

				runes := []rune(s)
				runes[0] = unicode.ToLower(runes[0])
				runes[len(runes)-1] = unicode.ToLower(runes[len(runes)-1])

				return string(runes), nil
			},
		}...)

		analyser.Analyse(text)
		topWords := analyser.Top(10)

		require.Equal(t, expectedTop, topWords)
	})
}
