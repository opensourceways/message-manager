/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller file scan service.
package controller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"github.com/opensourceways/message-manager/common/domain/allerror"
)

// TestControllerHttpError the unit test for the function
func TestControllerHttpError(t *testing.T) {
	convey.Convey("test func httpErro, err nil", t, func() {
		errcode, errString := httpError(nil)
		convey.So(errcode, convey.ShouldEqual, http.StatusOK)
		convey.So(errString, convey.ShouldEqual, "")
	})

	convey.Convey("test func httpErro, err branch test", t, func() {
		err := fmt.Errorf("default err")
		errcode1, errString1 := httpError(err)
		convey.So(errcode1, convey.ShouldEqual, http.StatusInternalServerError)
		convey.So(errString1, convey.ShouldEqual, "system_error")

		newErr1 := allerror.NewNotFound(allerror.ErrorCodeModelNotFound, "model not found")
		errcode2, errString2 := httpError(newErr1)
		convey.So(errcode2, convey.ShouldEqual, http.StatusNotFound)
		convey.So(errString2, convey.ShouldEqual, allerror.ErrorCodeModelNotFound)

		newErr2 := allerror.NewNoPermission("no permission")
		errcode3, errString3 := httpError(newErr2)
		convey.So(errcode3, convey.ShouldEqual, http.StatusForbidden)
		convey.So(errString3, convey.ShouldEqual, "no_permission")

		newErr3 := allerror.New(allerror.ErrorCodeAccessTokenInvalid, "error")
		errcode4, errString4 := httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeAccessTokenInvalid)

		newErr3 = allerror.New(allerror.ErrorCodeSessionIdMissing, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeSessionIdMissing)

		newErr3 = allerror.New(allerror.ErrorCodeSessionIdInvalid, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeSessionIdInvalid)

		newErr3 = allerror.New(allerror.ErrorCodeSessionNotFound, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeSessionNotFound)

		newErr3 = allerror.New(allerror.ErrorCodeSessionInvalid, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeSessionInvalid)

		newErr3 = allerror.New(allerror.ErrorCodeCSRFTokenMissing, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeCSRFTokenMissing)

		newErr3 = allerror.New(allerror.ErrorCodeCSRFTokenInvalid, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeCSRFTokenInvalid)

		newErr3 = allerror.New(allerror.ErrorCodeCSRFTokenNotFound, "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusUnauthorized)
		convey.So(errString4, convey.ShouldEqual, allerror.ErrorCodeCSRFTokenNotFound)

		newErr3 = allerror.New("Default_ERR", "error")
		errcode4, errString4 = httpError(newErr3)
		convey.So(errcode4, convey.ShouldEqual, http.StatusBadRequest)
		convey.So(errString4, convey.ShouldEqual, "Default_ERR")
	})
}
