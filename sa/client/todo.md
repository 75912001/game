# 基本要求
    深度思考 
    分析背景
    明确问题
    评审当前想法
    查阅可参考资料
    完成任务
    给出详细的解决方案.
    注重逻辑清晰和步骤明确.
    不要直接给出代码实现, 而是给出如何实现的思路.
    注重性能和可维护性.
    如若有更好的方案, 请一并给出.并给出理由和对比.
    可以完全访问当前项目的代码和资源.
# 背景
    需要增加 real time combat. 就是在地图上直接战斗.
# 问题
# 当前想法
    地图上有设定刷怪点, 刷怪点对应的是怪物组配置.C:\game\sa\client\cfg\enemy.group.yaml
    地图创建的时候, 读取刷怪点配置, 并在对应位置生成刷怪点实体.
    怪物在刷怪点附近随机巡逻.
# 可参考资料
    宠物配置表:
        C:\game\sa\client\cfg\pet.yaml
    地图配置表:
        C:\game\sa\client\res\tiled\map\map.2000003.tmx
        <property name="enemyGroupID" type="int" value="1"/>  刷怪点对应的怪物组ID (怪物组id,在 C:\game\sa\client\cfg\enemy.group.yaml 中配置)
# 任务
    设计并实现地图上刷怪点和怪物巡逻功能.