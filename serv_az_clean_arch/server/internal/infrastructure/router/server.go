package router

/*
import (
	"net/http"
	"server/internal/apperrors"
	"server/internal/config"
	"server/internal/email"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	GetPost(string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	InitImages(string, string)
}

type router struct {
	mux *mux.Router
}

func (router *router) InitImages(path string, absolutePath string) {
	router.mux.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(absolutePath))))
}

func (router *router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *router) GetPost(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodGet, http.MethodPost)
}

func (router *router) Get(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodGet)
}

func (router *router) Post(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodPost)
}

type server struct {
	repoLibrary repositories.RepoLibrary
	repoWords   repositories.RepoWords
	repoUsers   repositories.RepoUsers
	repoBackUp  repositories.BackUpCopyRepo
	router      Router
	sender      email.Sender
	logger      *logrus.Logger
	config      *config.Config
	blacklist   *blacklist
	tmpls       map[string]*template.Template
}

func NewServer(repoLibrary repositories.RepoLibrary, repoWords repositories.RepoWords, repoUsers repositories.RepoUsers, repoBackUp repositories.BackUpCopyRepo,
	logger *logrus.Logger, sender email.Sender, config *config.Config) *server {
	return &server{
		repoLibrary: repoLibrary,
		repoWords:   repoWords,
		repoUsers:   repoUsers,
		repoBackUp:  repoBackUp,
		router:      &router{mux: mux.NewRouter()},
		logger:      logger, config: config,
		sender: sender,
	}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHttp(w, r)
}

func (srv *server) initializeRoutes() {
	srv.initializeTemplates()
	srv.logger.Info("server INIT")
	srv.router.InitImages("/images/", "templates/images")
	//------------HOME----translate
	srv.router.Get("/", srv.contextExpire(srv.homeHandler()))
	srv.router.GetPost("/translate", srv.contextExpire(srv.getTranslationHandler()))
	srv.router.GetPost("/quick-answer", srv.contextExpire(srv.quickAnswerHandler()))
	//---------------user-CRUD---------------------------
	srv.router.GetPost("/registration", srv.contextExpire(srv.createUserHandler()))
	srv.router.GetPost("/login", srv.contextExpire(srv.loginHandler()))
	blackList := newBlacklist()
	srv.blacklist = blackList ///user-update-password
	srv.router.GetPost("/user-update", srv.jwtAuthentication(srv.updateUserHandler()))
	srv.router.GetPost("/user-update-password", srv.jwtAuthentication(srv.updateUserPasswordHandler()))
	srv.router.GetPost("/user-restore-password", srv.contextExpire(srv.restoreUserPasswordHandler()))
	srv.router.GetPost("/logout", srv.contextExpire(srv.logoutHandler()))
	srv.router.Get("/user-info", srv.jwtAuthentication(srv.getUserByIdHandler()))
	//-------TESTS---LEARN--------------------
	srv.router.GetPost("/test", srv.jwtAuthentication(srv.testHandler()))
	srv.router.GetPost("/learn", srv.jwtAuthentication(srv.learnHandler()))
	//-------TESTS--------thematic test----------------------
	srv.router.GetPost("/thematic/{theme:[0-9a-zA-Z_-]+}", srv.jwtAuthentication(srv.testUniversalHandler()))
	srv.router.GetPost("/test-thematic", srv.jwtAuthentication(srv.themesHandler()))
	//-------Update LIBRARY-------------------
	srv.router.GetPost("/info-users", srv.jwtAuthentication(srv.getAllUsersHandler()))
	srv.router.GetPost("/library-update", srv.jwtAuthentication(srv.updateLibraryHandler()))
	srv.router.Get("/library-download", srv.jwtAuthentication(srv.downloadHandler()))
}

func Run() {

	/*
		err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
*/
/*}

// "update_lib.html"
*/
