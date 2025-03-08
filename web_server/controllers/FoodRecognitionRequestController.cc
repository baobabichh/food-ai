#include "FoodRecognitionRequestController.h"

void FoodRecognitionRequestController::createRequest(const HttpRequestPtr &req, std::function<void(const HttpResponsePtr &)> &&callback) const
{
    const auto user_identity = getUserIdentityByUUID(req);
    if(!user_identity.isCorrect())
    {
        responseWithError(callback, "Incorrect UUID.");
        return;
    }

    const std::string& img_base64 = req->getParameter("img_base64");
    const std::string& img_type = req->getParameter("img_type");

    if(img_type.size() > 1 && img_base64.size() > 1) {}
    else 
    {
        responseWithError(callback, "Incorrect img_type or img_base64.");
        return;
    }

    auto client = drogon::app().getDbClient("food_ai");
    try
    {
        {
            static const std::string query = "insert into FoodRecognitionRequests (UserID, ImgBase64, ImgType, Status) values (?, ?, ?, ?)";
            const auto result = client->execSqlSync(query, std::to_string(user_identity.id), img_base64, img_type, "1");
        }

        static const std::string query = "SELECT LAST_INSERT_ID()";
        const auto result = client->execSqlSync(query);
        size_t request_id = result[0][0].as<size_t>(); // Retrieve the ID
        if(request_id == 0)
        {
            responseWithError(callback, "Failed to create request_id.");
            return;
        }

        responseWithSuccess(callback, "Request created.", {{"RequestID", std::to_string(request_id)}});
        return;
    }
    catch(const drogon::orm::DrogonDbException &e)
    {
        responseWithError(callback, "Internal server error.", {{"Exception", e.base().what()}});
        return;
    }
}

void FoodRecognitionRequestController::checkRequestStatus(const HttpRequestPtr &req, std::function<void(const HttpResponsePtr &)> &&callback) const
{
    const auto user_identity = getUserIdentityByUUID(req);
    if(!user_identity.isCorrect())
    {
        responseWithError(callback, "Incorrect UUID.");
        return;
    }

    const std::string& request_id = req->getParameter("request_id");
    if(stringToSizeT(request_id) > 0) {}
    else 
    {
        responseWithError(callback, "Incorrect request_id.");
        return;
    }

    auto client = drogon::app().getDbClient("food_ai");
    try
    {
        static const std::string query = "select Status from FoodRecognitionRequests where UserID = ? and ID = ?";
        const auto result = client->execSqlSync(query, std::to_string(user_identity.id), request_id);
        if(result.empty())
        {
            responseWithError(callback, "There is no such request.");
            return;
        }
        size_t status = result[0]["Status"].as<size_t>();
        if(status == 0)
        {
            responseWithError(callback, "Incorrect status.");
            return;
        }

        responseWithSuccess(callback, "Request status.", {{"RequestStatus", int_to_status[status]}});
        return;
    }
    catch(const drogon::orm::DrogonDbException &e)
    {
        responseWithError(callback, "Internal server error.", {{"Exception", e.base().what()}});
        return;
    }
}

void FoodRecognitionRequestController::getRequest(const HttpRequestPtr &req, std::function<void(const HttpResponsePtr &)> &&callback) const
{
    const auto user_identity = getUserIdentityByUUID(req);
    if(!user_identity.isCorrect())
    {
        responseWithError(callback, "Incorrect UUID.");
        return;
    }

    const std::string& request_id = req->getParameter("request_id");
    if(stringToSizeT(request_id) >= 0) {}
    else 
    {
        responseWithError(callback, "Incorrect request_id.");
        return;
    }

    auto client = drogon::app().getDbClient("food_ai");
    try
    {
        static const std::string query = "select ID, ImgBase64, ImgType, Status, Response, CreateTS from FoodRecognitionRequests where UserID = ? and ID = ?";
        const auto result = client->execSqlSync(query, std::to_string(user_identity.id), request_id);
        if(result.empty())
        {
            responseWithError(callback, "There is no such request.");
            return;
        }

        auto row = result[0];

        responseWithSuccess(callback, "Request status.", 
            {
                {"ID", row["ID"].as<std::string>()},
                {"ImgBase64", row["ImgBase64"].as<std::string>()},
                {"ImgType", row["ImgType"].as<std::string>()},
                {"Status", int_to_status[stringToSizeT(row["Status"].as<std::string>())]},
                {"Response", row["Response"].as<std::string>()},
                {"CreateTS", row["CreateTS"].as<std::string>()},
            });
        return;
    }
    catch(const drogon::orm::DrogonDbException &e)
    {
        responseWithError(callback, "Internal server error.", {{"Exception", e.base().what()}});
        return;
    }
}
