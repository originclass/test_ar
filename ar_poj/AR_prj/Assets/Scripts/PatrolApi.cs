using System;
using System.Collections;
using System.Text;
using UnityEngine;
using UnityEngine.Networking;

namespace PatrolAR
{
    /// <summary>
    /// 巡检问题上报 API 客户端(对应 go-admin 后端接口)。
    /// 纯 C# 实现,不依赖任何第三方 JSON 库;字段与后端 DTO 保持一致(camelCase)。
    ///
    /// 接口:POST {BaseUrl}/admin-api/v1/app/patrol/issue/report  (免登录,Unity 直接上报)
    /// 成功:{"code":200,"msg":"操作成功","data":1}          data = 后端生成的问题 ID
    /// 失败:{"code":31000,"msg":"问题标题不能为空","data":null}
    /// </summary>
    public static class PatrolApi
    {
        /// <summary>改成你电脑的局域网 IP(手机必须与电脑同一 WiFi)。cmd 里 ipconfig 查看 IPv4 地址。</summary>
        public static string BaseUrl = "http://192.168.10.26:8888";

        [Serializable]
        private class PositionDto
        {
            public double x;
            public double y;
            public double z;
        }

        [Serializable]
        private class ReportBody
        {
            public string title;
            public string description;
            public string priority;
            public PositionDto position;
        }

        [Serializable]
        private class ReportResponse
        {
            public int code;
            public string msg;
            public long? data;
        }

        /// <summary>
        /// 上报一条巡检问题。完成后回调 result:
        ///   success=true  时 id 为后端生成的问题 ID;
        ///   success=false 时 error 为用户可读的失败原因。
        /// </summary>
        public static IEnumerator ReportIssueAsync(string title, string description, string priority,
            Vector3 position, Action<bool, string, long> onResult)
        {
            var body = new ReportBody
            {
                title = title,
                description = description ?? "",
                priority = priority,
                position = new PositionDto
                {
                    x = Round(position.x),
                    y = Round(position.y),
                    z = Round(position.z)
                }
            };

            string json = JsonUtility.ToJson(body);
            string url = BaseUrl.TrimEnd('/') + "/admin-api/v1/app/patrol/issue/report";

            using (var request = new UnityWebRequest(url, "POST"))
            {
                byte[] bodyRaw = Encoding.UTF8.GetBytes(json);
                request.uploadHandler = new UploadHandlerRaw(bodyRaw);
                request.downloadHandler = new DownloadHandlerBuffer();
                request.SetRequestHeader("Content-Type", "application/json");
                request.timeout = 15;

                yield return request.SendWebRequest();

                // 网络层失败(连不上 / 超时 / 5xx)
                if (request.result != UnityWebRequest.Result.Success)
                {
                    string reason = DescribeNetworkError(request);
                    onResult?.Invoke(false, "上报失败:" + reason, 0);
                    yield break;
                }

                // 应用层:解析后端返回的 code/msg
                var resp = ParseResponse(request.downloadHandler.text);
                if (resp == null)
                {
                    onResult?.Invoke(false, "上报失败:服务器返回格式异常(" + request.responseCode + ")", 0);
                    yield break;
                }

                if (resp.code == 200)
                {
                    onResult?.Invoke(true, resp.msg, resp.data ?? 0);
                }
                else
                {
                    // 例如 31000 问题标题不能为空 / 31001 优先级必须是 low/medium/high 之一
                    onResult?.Invoke(false, "上报失败:" + resp.msg + "(错误码 " + resp.code + ")", 0);
                }
            }
        }

        private static double Round(double v)
        {
            return Math.Round(v, 4);
        }

        private static ReportResponse ParseResponse(string text)
        {
            if (string.IsNullOrEmpty(text))
            {
                return null;
            }
            try
            {
                return JsonUtility.FromJson<ReportResponse>(text);
            }
            catch (Exception)
            {
                return null;
            }
        }

        private static string DescribeNetworkError(UnityWebRequest request)
        {
            switch (request.result)
            {
                case UnityWebRequest.Result.ConnectionError:
                    return "无法连接服务器,请检查:① 电脑后端是否启动 ② 手机与电脑是否同一 WiFi ③ 地址是否为电脑局域网 IP(" + request.error + ")";
                case UnityWebRequest.Result.DataProcessingError:
                    return "服务器响应处理失败(" + request.error + ")";
                case UnityWebRequest.Result.ProtocolError:
                    return "服务器返回错误码 HTTP " + request.responseCode;
                default:
                    return request.error ?? "未知错误";
            }
        }
    }
}
