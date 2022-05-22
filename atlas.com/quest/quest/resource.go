package quest

import (
	"atlas-quest/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	getQuest = "get_quest"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, _ *gorm.DB) {
	r := router.PathPrefix("/quests").Subrouter()
	//r.HandleFunc("/", registerClearCache(l)).Methods(http.MethodDelete)
	//r.HandleFunc("/", registerGetQuestByInfoNumber(l)).Methods(http.MethodGet).Queries("infoNumber", "{infoNumber}", "filter[search]", "{filter}")
	//r.HandleFunc("/{id}", registerGetQuestCheckEnd(l)).Methods(http.MethodGet).Queries("checkEnd", "{checkEnd}")
	r.HandleFunc("/{id}", registerGetQuest(l)).Methods(http.MethodGet)
	//r.HandleFunc("/{id}/infoNumber", registerGetQuestInfoNumber(l)).Methods(http.MethodGet).Queries("status", "{status}")
	//r.HandleFunc("/{id}/infoEx", registerGetQuestInfoNumberEx(l)).Methods(http.MethodGet).Queries("status", "{status}", "index", "{index}")
	//r.HandleFunc("/{id}/items/{itemId}", registerGetQuestItemInformation(l)).Methods(http.MethodGet)
	//r.HandleFunc("/{id}", registerClearQuestCache(l)).Methods(http.MethodDelete)
	//r.HandleFunc("/items/skillBooks", registerGetSkillBooksFromQuests(l)).Methods(http.MethodGet)
}

type IdHandler func(questId uint32) http.HandlerFunc

func ParseId(l logrus.FieldLogger, next IdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		questId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse questId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(questId))(w, r)
	}
}

func registerGetQuest(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getQuest, func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(questId uint32) http.HandlerFunc {
			return handleGetQuest(l)(span)(questId)
		})
	})
}

func handleGetQuest(l logrus.FieldLogger) func(span opentracing.Span) func(questId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(questId uint32) http.HandlerFunc {
		return func(questId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, _ *http.Request) {
				_, err := GetById(l)(questId)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
