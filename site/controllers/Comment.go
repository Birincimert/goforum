package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"goforum/site/helpers"
	"goforum/site/models"

	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
)

type SiteComments struct {
	Store *sessions.CookieStore
}

// CommentAdd, yeni yorum veya yanıt ekler ve kullanıcıyı post sayfasına geri yönlendirir
func (c SiteComments) CommentAdd(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form verisi okunamadı", http.StatusBadRequest)
		return
	}

	// Giriş kontrolü
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postIDStr := r.FormValue("post-id")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}

	var parentCommentID *uint
	if parentStr := r.FormValue("parent-comment-id"); parentStr != "" {
		if pID, err := strconv.ParseUint(parentStr, 10, 32); err == nil {
			uID := uint(pID)
			parentCommentID = &uID
		}
	}

	comment := models.Comment{
		UserID:          user.ID,
		Name:            user.Username,
		Content:         r.FormValue("comment-content"),
		PostID:          uint(postID),
		ParentCommentID: parentCommentID,
	}

	if err := comment.Add(); err != nil {
		http.Error(w, "Yorum eklenirken hata oluştu", http.StatusInternalServerError)
		return
	}

	// Post'u ID ile al
	post := models.Post{}.Get(postID)
	if post.ID == 0 {
		http.Error(w, "Post bulunamadı", http.StatusNotFound)
		return
	}

	// Slug ile detay sayfasına yönlendir
	http.Redirect(w, r, fmt.Sprintf("/post/%s#comments", post.Slug), http.StatusSeeOther)
	return
}

// CommentUpvote yorumu beğenir
func (c SiteComments) CommentUpvote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// login zorunlu
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Error(w, "Giriş gerekli", http.StatusUnauthorized)
		return
	}
	commentIDStr := ps.ByName("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz yorum ID", http.StatusBadRequest)
		return
	}

	if err := (models.CommentVote{}).SetVote(user.ID, uint(commentID), 1); err != nil {
		http.Error(w, "Beğeni eklenirken hata oluştu", http.StatusInternalServerError)
		return
	}
	count, err := (models.CommentVote{}).CountVotes(uint(commentID))
	if err != nil {
		http.Error(w, "Beğeni sayısı hesaplanamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	liked, _ := (models.CommentVote{}).IsLikedBy(user.ID, uint(commentID))
	_, _ = fmt.Fprintf(w, `{"success": true, "likes": %d, "liked": %t}`, count, liked)
}

// CommentLikeCount mevcut beğeni sayısını döndürür
func (c SiteComments) CommentLikeCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_ = r // unused param
	commentIDStr := ps.ByName("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz yorum ID", http.StatusBadRequest)
		return
	}
	count, err := (models.CommentVote{}).CountVotes(uint(commentID))
	if err != nil {
		http.Error(w, "Beğeni sayısı alınamadı", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"success": true, "likes": %d}`, count)
}

// CommentIsLiked mevcut kullanıcı bu yorumu beğenmiş mi kontrol eder
func (c SiteComments) CommentIsLiked(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Error(w, "Giriş gerekli", http.StatusUnauthorized)
		return
	}
	commentIDStr := ps.ByName("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz yorum ID", http.StatusBadRequest)
		return
	}

	liked, err := (models.CommentVote{}).IsLikedBy(user.ID, uint(commentID))
	if err != nil {
		http.Error(w, "Durum alınamadı", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"success": true, "liked": %t}`, liked)
}

// DeleteOwnComment: kullanıcı yalnızca kendi yorumunu silebilir
func (c SiteComments) DeleteOwnComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, ok := helpers.GetCurrentUser(r, c.Store)
	if !ok {
		http.Redirect(w, r, "/login?return_url=/profile", http.StatusSeeOther)
		return
	}
	cidStr := ps.ByName("id")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil || cid == 0 {
		http.Error(w, "Geçersiz yorum ID", http.StatusBadRequest)
		return
	}

	cm, err := (models.Comment{}).GetByID(uint(cid))
	if err != nil || cm.ID == 0 {
		http.Error(w, "Yorum bulunamadı", http.StatusNotFound)
		return
	}
	if cm.UserID != user.ID {
		_ = helpers.SetAlert(w, r, "Bu yorumu silme yetkiniz yok.", c.Store)
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	if err := (models.Comment{}).Delete(cm.ID); err != nil {
		http.Error(w, "Yorum silinirken hata oluştu", http.StatusInternalServerError)
		return
	}

	_ = helpers.SetAlert(w, r, "Yorum silindi.", c.Store)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
