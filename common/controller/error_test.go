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

		newErr2 := allerror.NewNoPermission("no permission")
		errcode3, errString3 := httpError(newErr2)
		convey.So(errcode3, convey.ShouldEqual, http.StatusForbidden)
		convey.So(errString3, convey.ShouldEqual, "no_permission")
	})
}
