# **Manufacturing Reference Manual**

This manual explains how factories produce goods in three phases: **Allocation**, **Production**, and **Delivery**. It covers the key terms, how resources are used, and how factory units operate.

---

## **Key Concepts**

### **Factory Unit**
A **factory unit** is part of a colony or ship that builds goods. Each unit:
- Uses resources like **fuel**, **materials**, and **labor**.
- Tracks progress using a **backlog** (unfinished and finished goods).
- Produces goods over several turns.

---

### **Resources**
Factories need three types of resources:

- **Fuel** – Powers the factory.
- **Materials** – Split into:
    - **METs** (metallic materials)
    - **NMTs** (non-metallic materials)
- **Labor** – Includes:
    - **Professionals** (engineers, skilled workers)
    - **Unskilled Workers**
    - **Automatons** (robotic labor; substitutes for unskilled workers)

Factories can use up to **20 material units** (METs and NMTs combined) per turn.

---

### **Backlog**
Each factory unit has a **backlog**, which holds:

- **Raw Goods**: Items in progress, tracked by how complete they are (25%, 50%, 75%).
- **Finished Goods**: Items that are 100% complete and ready for delivery.

The backlog is part of the **colony or ship’s inventory**.

---

## **Manufacturing Phases**

---

### **1. Allocation**

This is the first phase of manufacturing.

- The **allocator** checks how many resources the factory unit can use.
- It asks the unit how much it takes to build one item (fuel, METs, NMTs, labor).
- It also asks:
    1. Fuel required per item
    2. METs required per item
    3. NMTs required per item
    4. Professionals required per item
    5. Unskilled labor required per item
    6. Maximum number of items the unit can produce
    7. Number of items already in the backlog

> ⚙️ **Important:** Automatons are treated as unskilled labor. The allocator prefers to assign automatons first.

- The allocator compares the unit’s needs to available resources and assigns fuel, labor, and materials.
- The factory unit always gets **fuel and labor for the full number of items**.
- It only gets **materials for the number of new items**, since items already in progress don’t need more materials.

**Example**:  
If the factory is told to build 10 items and already has 3 started, it gets:
- Fuel and labor for 10 items
- METs/NMTs for only **7 new items**

---

### **2. Production**

This is where work gets done.

The factory uses its assigned resources in this order:

1. **Finish raw goods** (75% → 100%)
2. **Improve raw goods** (50% → 75%)
3. **Improve raw goods** (25% → 50%)
4. **Create new raw goods** (adds items at 25% complete)

> 🔧 You can’t start a task unless you have **all the resources** needed to complete that stage.

- Labor (professionals + unskilled/automatons) is used in every step.
- METs and NMTs are only used when creating **new** raw goods.

---

### **3. Delivery**

This is the final step.

- The factory transfers **finished goods (100% complete)** from the backlog to the **ship or colony’s inventory**.
- **Goods are never delivered in an “assembled” state.**

> 🚫 If the ship or colony is full or damaged in battle, finished goods may not be delivered. The game engine handles these situations.

---

## **Summary**

| Phase        | What Happens                                                                 |
|--------------|------------------------------------------------------------------------------|
| Allocation   | Assigns fuel, labor, and materials based on factory needs and constraints   |
| Production   | Uses resources to finish, improve, or start items in the backlog            |
| Delivery     | Moves finished goods into the ship or colony inventory                      |

