import multiprocessing
import pandas as pd
import bcrypt

def hash_password(password, salt):
    hashed_password = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed_password.decode('utf-8')

def process_row(row):
    salt = bcrypt.gensalt()
    updated_row = dict(row)
    print(row['Password'])
    updated_row['Password_Hash'] = hash_password(row['Password'], salt)
    return updated_row

def main():
    pool = multiprocessing.Pool()


    df = pd.read_excel('C:\\Project\\Dealer.xlsx')
    processed_rows = pool.map(process_row, df.to_dict('records'))
    processed_df = pd.DataFrame(processed_rows)
    processed_df.to_excel('C:\\Project\\Dealer.xlsx', index=False)

if __name__ == '__main__':
    main()