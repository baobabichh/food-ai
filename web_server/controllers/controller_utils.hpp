#include "nlohmann/json.hpp"
#include <drogon/orm/Exception.h>
#include <drogon/drogon.h>
#include <string>

using namespace drogon;

inline const std::string createStatusJson(const std::string& type, const std::string& message, const std::vector<std::pair<std::string, std::string>>& add_info = {})
{
    nlohmann::json object{};
    object["Status"] = type;
    object["Message"] = message;
    for(const auto& [key, value] : add_info)
    {
        if(key == "Status" || key == "Message")
        {
            std::cerr << "Invalid add_info key: " << key << " value: " << value << "\n";
            continue;
        }
        object[key] = value;
    }

    return object.dump();
}

inline void responseWithError(std::function<void (const HttpResponsePtr &)> &callback, const std::string& message, const std::vector<std::pair<std::string, std::string>>& add_info = {})
{
    auto response = HttpResponse::newHttpResponse();
    response->setStatusCode(HttpStatusCode::k200OK);
    response->setBody(createStatusJson("Error", message, add_info));
    response->setContentTypeCode(drogon::CT_APPLICATION_JSON);
    callback(response);
}

inline void responseWithSuccess(std::function<void (const HttpResponsePtr &)> &callback, const std::string& message, const std::vector<std::pair<std::string, std::string>>& add_info = {})
{
    auto response = HttpResponse::newHttpResponse();
    response->setStatusCode(HttpStatusCode::k200OK);
    response->setBody(createStatusJson("Success", message, add_info));
    response->setContentTypeCode(drogon::CT_APPLICATION_JSON);
    callback(response);
}

struct UserIdentity
{
    size_t id{0};
    std::string uuid{};

    constexpr bool isCorrect()const
    {
        return id != 0 && uuid.size();
    }
};

inline UserIdentity getUserIdentityByUUID(const HttpRequestPtr &req)
{
    UserIdentity user_identity{};

    const std::string& uuid = req->getParameter("uuid");
    if(uuid.empty())
    {
        return UserIdentity{};
    }

    auto client = drogon::app().getDbClient("food_ai");
    try
    {
        static const std::string query = "select ID from Users where UUID = ?";
        const auto result = client->execSqlSync(query, uuid);
        if(result.empty())
        {
            return UserIdentity{};
        }



        user_identity.id = result[0]["ID"].as<size_t>();
        user_identity.uuid = uuid;

        return user_identity;
    }
    catch(const drogon::orm::DrogonDbException &e)
    {
        return UserIdentity{};
    }
}


inline size_t stringToSizeT(const std::string& str)
{
    try
    {
        return std::stoull(str);
    }
    catch(...)
    {
        return 0;
    }
}