import emptyImage from "../assets/empty.png";

export default function EmptyTasksRow() {
  return (
    <div className="flex flex-col items-center justify-center gap-4 px-6 py-12 text-center">
      <img className="h-48 w-auto max-w-full" src={emptyImage} alt="No tasks" />
      <div className="space-y-2">
        <p className="text-lg font-semibold text-gray-700">Chưa có công việc nào</p>
        <p className="text-sm text-gray-500">Bấm "Add Task" để tạo công việc mới.</p>
      </div>
    </div>
  );
}
