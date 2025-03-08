#pragma once

#include <drogon/HttpController.h>
#include "controllers/controller_utils.hpp"

using namespace drogon;

class UserController : public drogon::HttpController<UserController>
{
  public:
    METHOD_LIST_BEGIN
    ADD_METHOD_TO(UserController::registerUser, "/registerUser", Get); 
    ADD_METHOD_TO(UserController::loginUser, "/loginUser", Get); 
    METHOD_LIST_END

    void registerUser(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback);
    void loginUser(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback);
};
