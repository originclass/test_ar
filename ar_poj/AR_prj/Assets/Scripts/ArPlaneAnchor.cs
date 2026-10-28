using UnityEngine;

namespace PatrolAR
{
    /// <summary>
    /// 把 AR 平面变成"可视化 + 可点击"的实体。
    /// 挂载在 AR Plane Manager 的 Plane Prefab 上(由 AR Foundation 在每个平面被发现时自动实例化)。
    /// 职责:
    ///   1. 平面出现时用淡蓝色渲染出来(题目要求 2:检测并可视化水平平面);
    ///   2. 提供 TryRaycastPlane(命中点) 供 TapToPlaceReporter 查询"是否点在了平面上"。
    /// </summary>
    public class ArPlaneAnchor : MonoBehaviour
    {
        /// <summary>平面可视化颜色(淡蓝半透明,可随意改)。</summary>
        private static readonly Color PlaneColor = new Color(0.25f, 0.55f, 1f, 0.35f);

        private MeshRenderer meshRenderer;

        private void Awake()
        {
            meshRenderer = GetComponent<MeshRenderer>();
        }

        private void Start()
        {
            ApplyVisual();
        }

        /// <summary>用半透明材质渲染平面(运行时克隆材质实例,不改动共享材质资产)。</summary>
        private void ApplyVisual()
        {
            if (meshRenderer == null)
            {
                return;
            }
            Material mat = meshRenderer.material; // 访问 .material 会自动克隆实例
            if (mat.HasProperty("_BaseColor"))
            {
                mat.SetColor("_BaseColor", PlaneColor); // URP Lit
            }
            mat.color = PlaneColor; // Built-in Standard
            meshRenderer.shadowCastingMode = UnityEngine.Rendering.ShadowCastingMode.Off;
            meshRenderer.receiveShadows = false;
        }

        /// <summary>
        /// 判断一条屏幕射线是否命中了这个平面(供 TapToPlaceReporter 调用)。
        /// 返回 true 时 hitPoint 为命中的世界坐标。
        /// </summary>
        public bool TryRaycastPlane(Ray ray, out Vector3 hitPoint)
        {
            hitPoint = default;
            var collider = GetComponent<Collider>();
            if (collider == null)
            {
                return false;
            }
            if (!collider.Raycast(ray, out RaycastHit hit, 10f))
            {
                return false;
            }
            hitPoint = hit.point;
            return true;
        }
    }
}
