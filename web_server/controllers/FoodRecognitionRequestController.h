#pragma once

#include <drogon/HttpController.h>
#include "controllers/controller_utils.hpp"

using namespace drogon;

class FoodRecognitionRequestController : public drogon::HttpController<FoodRecognitionRequestController>
{
  public:
    METHOD_LIST_BEGIN
    ADD_METHOD_TO(FoodRecognitionRequestController::createRequest, "/createRequest", Post);
    ADD_METHOD_TO(FoodRecognitionRequestController::checkRequestStatus, "/checkRequestStatus", Get);
    ADD_METHOD_TO(FoodRecognitionRequestController::getRequest, "/getRequest", Get);
    METHOD_LIST_END

    void createRequest(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback) const;
    void checkRequestStatus(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback) const;
    void getRequest(const HttpRequestPtr& req, std::function<void (const HttpResponsePtr &)> &&callback) const;

    inline static std::unordered_map<size_t, std::string> int_to_status
    {
        {1, "Pending"},
        {2, "Processing"},
        {3, "Error"},
        {4, "Success"},
    };
};
