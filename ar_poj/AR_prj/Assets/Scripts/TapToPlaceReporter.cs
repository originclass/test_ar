using System.Collections;
using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.XR.ARFoundation;
using UnityEngine.XR.ARSubsystems;

namespace PatrolAR
{
    /// <summary>
    /// 核心 AR 交互脚本(题目 3.1 的主体):
    ///   1. 点击屏幕 → 如果点在 UI 上,忽略(要求 6:按钮/输入框不误放标记);
    ///   2. AR Raycast 命中检测到的水平平面 → 在命中位置生成标记(要求 2、3);
    ///   3. 弹出信息面板填写标题/描述/优先级 → 提交到 go-admin 后端(要求 4、5);
    ///   4. 成功/失败用 Toast 明确提示(要求 5)。
    /// 挂载到一个空物体(例如 "ArInteraction"),同场景放置 AR Session / AR Plane Manager / AR Raycast Manager。
    /// </summary>
    public class TapToPlaceReporter : MonoBehaviour
    {
        [Header("AR 管理器(场景中必须存在)")]
        public ARRaycastManager raycastManager;
        public ARPlaneManager planeManager;

        [Header("标记设置")]
        [Tooltip("标记预制体(默认用内置立方体,无需自己建)")]
        public GameObject markerPrefab;
        [Tooltip("标记缩放")]
        public float markerScale = 0.05f;

        [Header("UI 引用(在 Inspector 里拖入)")]
        [Tooltip("状态 Toast 文本(ReportCanvas 下的 StatusText)")]
        public UnityEngine.UI.Text statusText;
        [Tooltip("填写面板(ReportPanel,默认隐藏)")]
        public GameObject reportPanel;
        [Tooltip("标题输入框(必填)")]
        public UnityEngine.UI.InputField titleInput;
        [Tooltip("描述输入框(可选)")]
        public UnityEngine.UI.InputField descriptionInput;
        [Tooltip("优先级下拉(low/medium/high)")]
        public UnityEngine.UI.Dropdown priorityDropdown;
        [Tooltip("填写面板里的 提交 按钮")]
        public UnityEngine.UI.Button submitButton;
        [Tooltip("填写面板里的 取消 按钮")]
        public UnityEngine.UI.Button cancelButton;
        [Tooltip("Toast 显示秒数")]
        public float toastSeconds = 3f;

        /// <summary>当前待提交的标记(填表期间保持显示,提交成功后变色存档)。</summary>
        private GameObject pendingMarker;
        private bool busy; // 防重复提交

        private static readonly XRRaycastHit[] emptyHits = new XRRaycastHit[0];

        /// <summary>默认引导语(无其他提示时显示)。</summary>
        private const string DefaultHint = "移动手机扫描地面,出现淡蓝色平面后,点击它放置标记";

        private void Start()
        {
            if (markerPrefab == null)
            {
                // 没拖 prefab 就用 Unity 内置立方体,保证开箱即用
                markerPrefab = GameObject.CreatePrimitive(PrimitiveType.Cube);
                markerPrefab.GetComponent<Collider>().enabled = false;
                markerPrefab.SetActive(false);
            }

            if (submitButton != null) submitButton.onClick.AddListener(OnSubmitClicked);
            if (cancelButton != null) cancelButton.onClick.AddListener(OnCancelClicked);
            if (reportPanel != null) reportPanel.SetActive(false);
#if ENABLE_INPUT_SYSTEM && !ENABLE_LEGACY_INPUT_MANAGER
            UnityEngine.InputSystem.EnhancedTouch.EnhancedTouchSupport.Enable();
#endif
            ShowHint(DefaultHint);
        }

        /// <summary>
        /// 当前帧是否有"刚刚按下"的触摸(双输入系统兼容):
        /// 新 Input System 启用时走 UnityEngine.InputSystem.EnhancedTouch,否则走旧 Input.GetTouch。
        /// </summary>
        private bool TryGetBeginTouch(out Vector2 screenPos)
        {
#if ENABLE_INPUT_SYSTEM && !ENABLE_LEGACY_INPUT_MANAGER
            // ---- 仅新 Input System ----
            var touches = UnityEngine.InputSystem.EnhancedTouch.Touch.activeTouches;
            for (int i = 0; i < touches.Count; i++)
            {
                if (touches[i].phase == UnityEngine.InputSystem.TouchPhase.Began)
                {
                    screenPos = touches[i].screenPosition;
                    return true;
                }
            }
            screenPos = default;
            return false;
#else
            // ---- 旧 Input Manager(或 Both 模式)----
            if (Input.touchCount > 0 && Input.GetTouch(0).phase == TouchPhase.Began)
            {
                screenPos = Input.GetTouch(0).position;
                return true;
            }
            screenPos = default;
            return false;
#endif
        }

        private void Update()
        {
            if (busy || reportPanel == null || reportPanel.activeSelf)
            {
                return; // 填表中/提交中不允许再放标记
            }

            Vector2 touchPos;
            if (!TryGetBeginTouch(out touchPos))
            {
                return;
            }

            // 要求 6:点在任意 UI(按钮/输入框/面板)上时不放置标记
            if (EventSystem.current != null && EventSystem.current.IsPointerOverGameObject())
            {
                return;
            }

            TryPlaceMarker(touchPos);
        }

        /// <summary>要求 2、3:AR 射线检测水平平面,命中则在命中点放置标记。</summary>
        private void TryPlaceMarker(Vector2 screenPos)
        {
            if (raycastManager == null)
            {
                ShowToast("未配置 ARRaycastManager", true);
                return;
            }

            var hits = new System.Collections.Generic.List<ARRaycastHit>();
            // TrackableType.PlaneWithinBounds:只认已检测到的平面区域(AR Foundation 6.x 用单数)
            if (raycastManager.Raycast(screenPos, hits, TrackableType.PlaneWithinBounds) && hits.Count > 0)
            {
                Pose hitPose = hits[0].pose;

                if (pendingMarker != null)
                {
                    Destroy(pendingMarker);
                }
                pendingMarker = Instantiate(markerPrefab, hitPose.position, Quaternion.identity);
                pendingMarker.transform.localScale = Vector3.one * markerScale;
                pendingMarker.SetActive(true);

                SetMarkerColor(pendingMarker, new Color(1f, 0.6f, 0f)); // 橙色=待提交

                if (reportPanel != null)
                {
                    reportPanel.SetActive(true);
                    if (titleInput != null) titleInput.ActivateInputField();
                }
                ShowHint("已放置标记:请填写问题信息,填完点提交");
            }
            else
            {
                ShowHint("未检测到平面:请对准地面/桌面,缓慢移动手机");
            }
        }

        /// <summary>提交按钮:校验 → 调 PatrolApi → 按结果提示(要求 4、5)。</summary>
        private void OnSubmitClicked()
        {
            if (busy) return;
            if (pendingMarker == null)
            {
                ShowToast("请先在平面上点击放置标记", true);
                return;
            }
            if (titleInput == null || string.IsNullOrWhiteSpace(titleInput.text))
            {
                ShowToast("标题不能为空", true);
                return;
            }

            string priority = priorityDropdown != null && priorityDropdown.options.Count > 0
                ? priorityDropdown.options[priorityDropdown.value].text
                : "low";

            busy = true;
            ShowToast("正在提交…");
            StartCoroutine(PatrolApi.ReportIssueAsync(
                titleInput.text.Trim(),
                descriptionInput != null ? descriptionInput.text.Trim() : "",
                priority,
                pendingMarker.transform.position,
                (success, message, id) =>
                {
                    busy = false;
                    if (success)
                    {
                        SetMarkerColor(pendingMarker, new Color(0.2f, 0.8f, 0.4f)); // 绿色=已提交
                        pendingMarker = null; // 存档,不再复用
                        ClearForm();
                        if (reportPanel != null) reportPanel.SetActive(false);
                        ShowToast("提交成功!问题编号 #" + id);
                    }
                    else
                    {
                        ShowToast(message, true);
                        // 标记保留在场景里,用户改完可直接重试
                    }
                }));
        }

        private void OnCancelClicked()
        {
            if (pendingMarker != null)
            {
                Destroy(pendingMarker);
                pendingMarker = null;
            }
            ClearForm();
            if (reportPanel != null) reportPanel.SetActive(false);
            ShowToast("已取消");
        }

        private void ClearForm()
        {
            if (titleInput != null) titleInput.text = "";
            if (descriptionInput != null) descriptionInput.text = "";
            if (priorityDropdown != null) priorityDropdown.value = 0;
        }

        private void SetMarkerColor(GameObject marker, Color color)
        {
            var renderer = marker.GetComponent<Renderer>();
            if (renderer != null)
            {
                Material mat = renderer.material; // 克隆实例,避免改到共享材质
                if (mat.HasProperty("_BaseColor")) mat.SetColor("_BaseColor", color); // URP
                mat.color = color;                                                    // Built-in
            }
        }

        /// <summary>一次性提示(几秒后消失,恢复引导语)。</summary>
        private void ShowToast(string msg, bool isError = false)
        {
            if (statusText == null) return;
            statusText.text = msg;
            statusText.color = isError ? new Color(1f, 0.35f, 0.35f) : Color.white;
            StopAllCoroutines();
            StartCoroutine(FadeToastAfter(toastSeconds));
        }

        /// <summary>常驻提示(不自动消失,直到被下一个提示替换)。</summary>
        private void ShowHint(string msg)
        {
            if (statusText == null) return;
            statusText.text = msg;
            statusText.color = Color.white;
            StopAllCoroutines();
        }

        private IEnumerator FadeToastAfter(float seconds)
        {
            yield return new WaitForSeconds(seconds);
            if (statusText != null)
            {
                statusText.text = DefaultHint; // Toast 结束后回到引导语,而不是空白
            }
        }
    }
}
