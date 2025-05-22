import sqlite3
import tkinter as tk
from PIL import Image, ImageTk

DB_PATH = "/home/mxz-11/Desktop/robotica/data_treatment/ddbb/winery.db"
TABLE_NAME = "winery_data"

def get_latest_data():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    cursor.execute(f"""
        SELECT co2, temperature, humidity, date
        FROM {TABLE_NAME}
        ORDER BY datetime(date) DESC
        LIMIT 1
    """)
    row = cursor.fetchone()
    conn.close()
    return row

def prepare_items(data):
    co2, temp, hum, date = data
    return [
        ("CO₂: {} ppm".format(co2), "co.png"),
        ("Temperature: {} °C".format(temp), "temp.png"),
        ("Humidity: {} %".format(hum), "hum.png"),
        ("Date: {}".format(date), "cal.png")
    ]

def show_next_message(event=None):
    global current_index
    current_index = (current_index + 1) % len(items)
    text, img_path = items[current_index]
    message_var.set(text)
    update_image(img_path)
    center_widgets()

def update_image(path):
    try:
        img = Image.open(path)
        img = img.resize((120, 120), Image.Resampling.LANCZOS)
        photo = ImageTk.PhotoImage(img)
        image_label.config(image=photo)
        image_label.image = photo
    except Exception as e:
        print(f"Error loading image {path}: {e}")

def center_widgets():
    text_label.update_idletasks()
    text_width = text_label.winfo_width()
    total_width = image_width + gap + text_width

    start_x = (screen_width - total_width) // 2
    image_label.place(x=start_x, y=center_y)
    text_label.place(x=start_x + image_width + gap, y=center_y)

def float_image(offset=0, direction=1):
    new_offset = offset + direction
    image_label.place_configure(y=center_y + new_offset)
    if abs(new_offset) >= 10:
        direction *= -1
    root.after(50, lambda: float_image(new_offset, direction))

def exit_app(event):
    root.destroy()

def poll_for_new_data():
    global items, current_index, last_data
    new_data = get_latest_data()
    if new_data and new_data != last_data:
        print("New data detected, updating display...")
        last_data = new_data
        items = prepare_items(new_data)
        current_index = 0
        message_var.set(items[current_index][0])
        update_image(items[current_index][1])
        center_widgets()
    root.after(5000, poll_for_new_data)  # Vuelve a comprobar en 5 segundos

# --- INTERFAZ ---
root = tk.Tk()
root.configure(bg="black")
root.attributes("-fullscreen", True)
root.bind("<Escape>", exit_app)

screen_width = root.winfo_screenwidth()
screen_height = root.winfo_screenheight()
center_y = screen_height // 2

image_width = 120
gap = 30

# Crear widgets
image_label = tk.Label(root, bg="black")
text_label = tk.Label(root, font=("Arial", 50),
                      bg="black", fg="white", justify="left", anchor="w")
message_var = tk.StringVar()
text_label.config(textvariable=message_var)

# Datos iniciales
last_data = get_latest_data()
if last_data:
    items = prepare_items(last_data)
else:
    items = [("No data available", "error.png")]

current_index = 0
message_var.set(items[current_index][0])
update_image(items[current_index][1])
center_widgets()

float_image()
poll_for_new_data()  # Empieza la comprobación periódica
root.bind("<Button-1>", show_next_message)
root.mainloop()
