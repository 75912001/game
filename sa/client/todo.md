# 基本要求
    深度思考, 分析背景和任务.
    给出详细的解决方案.
    注重逻辑清晰和步骤明确.
    不要直接给出代码实现, 而是给出如何实现的思路.
    注重性能和可维护性.
# 背景
    现在在场景中增加一个 tile , 作为大型树木, 如雨伞状. 底部较窄, 顶部较宽.高度较高.
    按照当前的角色和场景的绘制顺序(Y轴排序). 会出现角色走到树木的上方时, 树木会遮挡角色的问题.
    应该是树冠部分作为 overhead 层绘制, 树干部分作为 ground 层绘制. 
    这样树干部分可以和角色进行Y轴排序.(现有方案)
    树冠部分作为 overhead 层绘制, 始终在角色的上方.(想要的方案)

    想通过代码实现自动拆分树木 tile 为 ground 和 overhead 两部分进行绘制.
    在树木的 tile 属性中增加一个标记, 标记该 tile 需要拆分为 ground 和 overhead 两部分.
    (目前增加的属性名称为 overheadRatio, 表示树冠部分占整个 tile 的比例, 取值范围 0.0 ~ 1.0)

    通过代码读取该 tile 的属性, 自动拆分为 ground 和 overhead 两部分
        是否可以创建一个新的 缓存, 专门用于存放拆分后的 ground 和 overhead tile ?

    植物信息位置: C:\game\sa\client\res\plant
    地图信息位置: C:\game\sa\client\res\tiled\map
# 任务
    分析背景.
    如何通过代码实现自动拆分该 tile 为 ground 和 overhead 两部分进行绘制.
    将拥有 overheadRatio 属性的 tile 进行拆分.
        加载配置的时候进行拆分.
        绘制的时候进行绘制.
