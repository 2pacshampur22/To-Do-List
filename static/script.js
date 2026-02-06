const API_URL = "http://localhost:8080/tasks/";

const listElement = document.getElementById("taskList");
const inputElement = document.getElementById("taskInput");

document.addEventListener("DOMContentLoaded", () => {
    fetchTasks();
});
async function fetchTasks() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error("Ошибка сети");
        
        const tasks = await response.json();
        
        renderTasks(tasks);
    } catch (error) {
        console.error("Ошибка:", error);
        listElement.innerHTML = `<p class="text-red-400 text-center">Не удалось загрузить задачи :(</p>`;
    }
}


async function createTask() {
    const nameInput = document.getElementById("nameInput"); 
    const descInput = document.getElementById("taskInput"); 

    const nameText = nameInput.value.trim();
    const descText = descInput.value.trim();

    if (!nameText) return alert("Введите название задачи!");

    try {
        await fetch(API_URL, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ 
                name: nameText, 
                description: descText 
            }) 
        });

        nameInput.value = "";
        descInput.value = "";
        
        fetchTasks(); 
    } catch (error) {
        console.error("Ошибка при создании:", error);
    }
}

function renderTasks(tasks) {
    const listElement = document.getElementById("taskList");
    listElement.innerHTML = "";

    if (!tasks || tasks.length === 0) {
        listElement.innerHTML = `<p class="text-slate-500 text-center">Список пуст.</p>`;
        return;
    }

    tasks.forEach(task => {
        const opacityClass = task.is_done ? "opacity-50" : "opacity-100";
        const textDecoration = task.is_done ? "line-through text-slate-500" : "text-slate-200";
        const borderClass = task.is_done ? "border-slate-600" : "border-indigo-500";
        
        const html = `
            <li class="flex justify-between items-center bg-slate-700 p-4 rounded-lg border-l-4 ${borderClass} ${opacityClass} transition-all duration-300">
                
                <div 
                    onclick="toggleTask(${task.id})" 
                    class="flex-1 flex flex-col cursor-pointer select-none mr-4"
                >
                    <span class="font-bold text-lg ${textDecoration}">
                        ${task.name}
                    </span>
                    
                    <span class="text-sm text-slate-400 break-words">
                        ${task.description}
                    </span>
                </div>
                
                <button 
                    onclick="deleteTask(${task.id})"
                    class="text-red-400 hover:text-red-300 hover:bg-slate-600 p-2 rounded-full transition shrink-0">
                    ✕
                </button>
            </li>
        `;

        listElement.insertAdjacentHTML("beforeend", html);
    });
}

async function deleteTask(id) {
    if (!confirm("Удалить задачу?")) return;

    try {
        await fetch(API_URL + id, { method: "DELETE" });
        fetchTasks();
    } catch (error) {
        console.error("Ошибка при удалении:", error);
    }
}

async function toggleTask(id) {
    try {
        await fetch(API_URL + id, { method: "PUT" });
        fetchTasks();
    } catch (error) {
        console.error("Ошибка при обновлении:", error);
    }
}